# orderedWorkerPool

orderedWorkerPool simple worker pool that executes work asynchronously but
delivers results in order of dispatch.

### Usage:
```
func main() {
    tasksCh, resultCh := orderedWorkerPool.New(runtime.NumCPU() * 2)
    go dispatchTasks(tasksCh, taskCount)
    for result := range resultCh {
        // Handle results
    }
}

func dispatchTasks(taskCh chan<- *Task) {
    for id := 0; id < 100; id++ {

        // This `Task` simply returns id after sleeping for 3ms
        args := []interface{}{id}
        f := func(args ...interface{}) interface{} {
            time.Sleep(time.Millisecond * 3)
            return args[0]
        }
        taskCh <- NewTask(id, args, f)
    }
    close(taskCh)
}
```

```
go get github.com/stringcode86/orderedWorkerPool
```