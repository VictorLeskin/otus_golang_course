package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func createStage(in In, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer close(out)

		// Запускаем оригинальный stage
		stageOut := stage(in)

		defer func() {
			for v := range stageOut {
				_ = v
			}
		}()

		for {
			select {
			case <-done:
				// Прерываем выполнение по сигналу
				return
			case val, ok := <-stageOut:
				if !ok {
					// Канал stage закрыт - завершаем
					return
				}

				// Пытаемся отправить значение с проверкой done
				select {
				case <-done:
					return
				case out <- val:
					// Продолжаем
				}
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		return nil
	}
	current := in
	for _, stage := range stages {
		// Создаем канал для текущего этапа с обработкой done
		current = createStage(current, done, stage)
	}

	return current
}
