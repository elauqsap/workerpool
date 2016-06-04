#### WorkerPool
My modifications to a blog [post](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/) from an employee at Malwarebytes

##### Goal
Package my modifications for easy reuse

##### Example
See example/example.go for the full code
``` go
// Initialize a Dispatcher
dispatcher := wp.NewDispatcher(2, 1024)
// Start the Dispatcher and create/register the Workers to the WorkerPool
dispatcher.Run()

// Queue two jobs for processing
dispatcher.WaitGroup.Add(2)
dispatcher.JobQueue <- SomeJob1{}
dispatcher.JobQueue <- SomeJob2{}

// Block main thread until processing in go routines completes
dispatcher.WaitGroup.Wait()
```
