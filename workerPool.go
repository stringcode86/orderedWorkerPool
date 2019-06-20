package orderedWorkerPool

import (
	"sort"
)

// New returns taskCh and results channels. Dispatch tasks on `taskCh` and
// receive results on `orderedResultCh`.
// Usage:
// ```
// func main() {
// 	 tasksCh, resultCh := orderedWorkerPool.New(runtime.NumCPU() * 2)
//	 go dispatchTasks(tasksCh, taskCount)
//	 for result := range resultCh {
//	 	// Handle results
//	 }
// }
//
// func dispatchTasks(taskCh chan<- *Task) {
// 	 for id := 0; id < 100; id++ {
//
//		// This `Task` simply returns id after sleeping for 3ms
//		args := []interface{}{id}
//		f := func(args ...interface{}) interface{} {
//			time.Sleep(time.Millisecond * 3)
//			return args[0]
//		}
//		taskCh <- NewTask(id, args, f)
//	}
//	close(taskCh)
// }
// ```
func New(workerCnt int) (taskCh chan<- *Task, orderedResultCh <-chan *Result) {
	_taskCh := make(chan *Task, 10)
	resultCh := make(chan *Result, 10)
	orderedResCh := make(chan *Result, 10)
	doneCh := make(chan bool, workerCnt)
	go workerDoneHandler(doneCh, workerCnt, resultCh)
	go resultHandler(resultCh, orderedResCh)
	for i := 0; i < workerCnt; i++ {
		go worker(_taskCh, resultCh, doneCh)
	}
	return _taskCh, orderedResCh
}

// Task encapsulates work do be done. `f` called with `args` and dispatched on 
// on results channel based on `id` (essentially sort order)
type Task struct {
	id   int
	args []interface{}
	f    func(args ...interface{}) interface{}
}

// NewTask with `id` (sort order), `args` to be passed to `f`
func NewTask(id int, args []interface{}, f func(args ...interface{}) interface{}) *Task {
	return &Task{id, args, f}
}

// Result of `Task.f` call with `Task.id`
type Result struct {
	Id    int
	Value interface{}
}

// worker pops tasks from `taskCh`. Executes task, sends result to `resultCh`.
// Once `taskCh` is closed, sends true to `doneCh`
func worker(taskCh <-chan *Task, resultCh chan<- *Result, doneCh chan<- bool) {
	for t := range taskCh {
		resultCh <- &Result{
			Id:    t.id,
			Value: t.f(t.args...),
		}
	}
	doneCh <- true
}

// resultHandler handles ordering of result. Dispatches them to `ordResCh` based
// on `Task.id`
func resultHandler(resCh <-chan *Result, ordResCh chan<- *Result) {
	cursor := -1
	resCache := make([]*Result, 0)
	for result := range resCh {
		if (cursor + 1) == result.Id {
			cursor = result.Id
			ordResCh <- result
		} else {
			resCache = append(resCache, result)
			sort.Slice(resCache, func(i, j int) bool {
				return resCache[i].Id < resCache[j].Id
			})
		}
		for len(resCache) > 0 && resCache[0].Id == (cursor+1) {
			result := resCache[0]
			resCache = resCache[1:]
			cursor = result.Id
			ordResCh <- result
		}
	}
	close(ordResCh)
}

// workerDoneHandler handles closing of unordered `resCh`
func workerDoneHandler(doneCh chan bool, wc int, resCh chan<- *Result) {
	doneCount := 0
	for range doneCh {
		doneCount++
		if doneCount == wc {
			close(doneCh)
			close(resCh)
		}
	}
}
