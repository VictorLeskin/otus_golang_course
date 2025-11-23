package hw05parallelexecution

import (
	"errors"
	"fmt"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

type Worker struct {
}

type WorkerPool struct {
	tasks          []Task
	processesCount int // count of processes
	maxErrors      int // max possible error

	chanTasks   chan Task
	chanResults chan error

	numWorkers int
	wg         sync.WaitGroup
	done       chan struct{}
	errorCount int
}

func (t *Worker) Run(task Task) {
	_ = task
}

func NewWorkerPool(tasks []Task, n, m int) *WorkerPool {
	// Place your code here.
	workerPool := WorkerPool{tasks: tasks,
		processesCount: n, maxErrors: m,
		chanTasks:   make(chan Task, len(tasks)),
		chanResults: make(chan error, len(tasks)),
		numWorkers:  min(len(tasks), n),
		errorCount:  0,
	}

	return &workerPool
}

func (t *WorkerPool) Run() error {

	for range t.numWorkers {
		t.wg.Add(1)
		go func() {
			defer t.wg.Done()
			for {
				select {
				case <-t.done:
					// Контекст отменен, завершаем работу
					return
				case task, ok := <-t.chanTasks:
					if !ok {
						// Канал закрыт, заданий больше нет
						return
					}
					err := task()
					if err != nil {
						t.chanResults <- err
					}

				}
			}
		}()
	}

	// Монитор результатов
	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		for {
			result, ok := <-t.chanResults
			if !ok {
				return
			}

			if result != nil {
				t.errorCount++
				if t.errorCount >= t.maxErrors {
					// Контекст отменен, завершаем работу
					close(t.done)
					return
				}
			}
		}
	}()

	for _, task := range t.tasks {
		t.chanTasks <- task
	}
	close(t.chanTasks)

	t.wg.Wait()
	close(t.chanResults)

	return nil
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	// Place your code here.
	workerPool := NewWorkerPool(tasks, n, m)

	fmt.Print(workerPool)

	return workerPool.Run()
}
