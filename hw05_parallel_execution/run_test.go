package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

// test with 2 tasks.
func TestRun10Tasks_IgnoreErrors(t *testing.T) {
	defer goleak.VerifyNone(t)
	tasksCount := 10
	tasks := make([]Task, 0, tasksCount)
	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			return err
		})
	}

	err := Run(tasks, 4, 0)
	require.NoError(t, err)
}

// test with 2 tasks.
func TestRun2Tasks(t *testing.T) {
	defer goleak.VerifyNone(t)
	tasks := make([]Task, 0, 2)
	tasks = append(tasks, func() error {
		time.Sleep(time.Millisecond * 2)
		return nil
	})
	tasks = append(tasks, func() error {
		time.Sleep(time.Millisecond * 3)
		return nil
	})

	err := Run(tasks, 1, 1)
	require.NoError(t, err)
}

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, int32(tasksCount), runTasksCount, "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}

// test with 2 tasks.
func TestRunTasks_1(t *testing.T) {
	// Создаем новый генератор с fixed seed	for stabilty of tests
	source := rand.NewSource(99)
	r := rand.New(source)

	tests := []struct {
		taskWOECount   int // count of task without error
		taskWECount    int // count of task with error
		workersCount   int
		maxErrorsCount int
		noError        bool
	}{
		{10, 0, 4, 0, true},
		{10, 0, 4, 1, true},
		{10, 1, 4, 1, false},
		{1, 9, 4, 8, false},
		{1, 9, 4, 9, false},
		{1, 9, 4, 10, true},
		{1, 0, 4, 0, true},
		{0, 1, 4, 2, true},
		{0, 1, 4, 1, false},
		{1, 0, 1, 0, true},
		{0, 1, 1, 2, true},
		{0, 1, 1, 1, false},
	}

	for n, tc := range tests {
		_ = n
		func() {
			defer goleak.VerifyNone(t)
			taskCount := tc.taskWOECount + tc.taskWECount
			tasks := make([]Task, 0, taskCount)
			for i := 0; i < tc.taskWECount; i++ {
				err := fmt.Errorf("error from task %d", i)
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					return err
				})
			}
			for i := 0; i < tc.taskWOECount; i++ {
				tasks = append(tasks, func() error {
					time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
					return nil
				})
			}

			r.Shuffle(len(tasks), func(i, j int) {
				tasks[i], tasks[j] = tasks[j], tasks[i]
			})

			err := Run(tasks, tc.workersCount, tc.maxErrorsCount)
			if tc.noError {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		}()
	}
}
