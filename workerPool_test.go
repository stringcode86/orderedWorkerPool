package orderedWorkerPool

import (
	"runtime"
	"testing"
	"time"
)

// TestOrderedWorkerPool creates `Task`s that simply return index as result,
// checks weather results arrive in order.
func TestOrderedWorkerPool(t *testing.T) {
	taskCount := 100000
	results := make([]int, taskCount)

	tasksCh, resultCh := New(runtime.NumCPU() * 2)
	go dispatchTasks(tasksCh, taskCount)

	// Cache the results in order they are received
	i := 0
	for result := range resultCh {
		results[i] = result.Value.(int)
		i++
	}

	// Verify that results arrived in correct order.
	for i := 0; i < taskCount; i++ {
		if results[i] != i {
			t.Error("Results out of order", results)
		}
	}
}

// TestOrderedWorkerPool creates `Task`s that simply return index as result,
// checks weather results arrive in order.
func TestOrderedWorkerPoolSmallNumberOfTasks(t *testing.T) {
	taskCount := 5
	results := make([]int, taskCount)

	tasksCh, resultCh := New(runtime.NumCPU() * 2)
	go dispatchTasks(tasksCh, taskCount)

	// Cache the results in order they are received
	i := 0
	for result := range resultCh {
		results[i] = result.Value.(int)
		i++
	}

	// Verify that results arrived in correct order.
	for i := 0; i < taskCount; i++ {
		if results[i] != i {
			t.Error("Results out of order", results)
		}
	}
}

func dispatchTasks(taskCh chan<- *Task, taskCount int) {
	for i := 0; i < taskCount; i++ {
		args := []interface{}{i}
		f := func(args ...interface{}) interface{} {
			time.Sleep(time.Millisecond * 3)
			return args[0]
		}
		taskCh <- NewTask(i, args, f)
	}
	close(taskCh)
}
