package scheduler

import (
	"proj3/concurrent"
	"proj3/func_objs"
	"proj3/png"
	"proj3/iohandler"
	"math/rand"
	"strings"
	"sync"
)


// Stealing Scheduler is the context for the Work Stealing Scheduler
type StealingScheduler struct {
	// this is the context for scheduler
	Mutex *sync.Mutex
	Group *sync.WaitGroup
	Deqs *([]*concurrent.Dequeue)
	Random int
	Tasks *([]png.ImageTask)
	ThreadCount int
	Workers *([]stealingWorker)
	DoneTaskCount *int
}

//Creates a new Work Stealing context
func NewStealingScheduler(capacity int, config Config) *StealingScheduler {
	deqs := make([](*concurrent.Dequeue), capacity)

	for i := 0; i < config.ThreadCount; i++ {
		deq := concurrent.NewDEQueue()
		deqs[i] = deq
	}

	// when initial scheduler, read the data from txt
	dataDirs := strings.Split(config.DataDirs, "+")
	infoMaps := iohandler.ReadTxt(dataDirs)
	tasks := iohandler.CreateImageTasks(infoMaps, dataDirs)
	workers := make([]stealingWorker, capacity)
	doneTaskCount := 0

	return &StealingScheduler{Mutex: &sync.Mutex{}, Group: &sync.WaitGroup{}, Deqs: &deqs, Workers: &workers, 
		ThreadCount: config.ThreadCount, Tasks: &tasks, DoneTaskCount: &doneTaskCount}
}

//PushTask inserts a new task into a queue of one of the threads in the pool
//You need to insure each thread is assigned an equal amount of work
func (scheduler *StealingScheduler) PushTask(idx int, task func_objs.Runnable, taskArg interface{}) {
	(*scheduler.Deqs)[idx].PushBottom(task, taskArg)
}

// Run starts the run() method for each thread in the pool
func (scheduler *StealingScheduler) Run() {
	scheduler.Group.Add(scheduler.ThreadCount)
	for i := 0; i < scheduler.ThreadCount; i++ {
		go (*scheduler.Workers)[i].run(i, scheduler.Deqs, scheduler.DoneTaskCount, len(*scheduler.Tasks), scheduler.Group)
	}
	scheduler.Wait()
}

// Done is a way for the application to tell the scheduler no more tasks will
// need to be handled by the scheduler. This method should notify the workers in some way
func (scheduler *StealingScheduler) Done() {
}

// Wait() is called when a goroutine wants to wait until all of the stealing threads in the pool have executed. 
// You may use a sync.Waitgroup in this implementation but only one.
func (scheduler *StealingScheduler) Wait() {
	scheduler.Group.Wait()
}

type stealingWorker struct {
	deq    *concurrent.Dequeue
}

// Creates a new stealer worker (i.e., a thread in the pool of workers)
func NewStealingWorker() *stealingWorker {
	deq := concurrent.NewDEQueue()
	return &stealingWorker{deq: deq}
}

// Run is the main function being executed by a stealing worker.
func (worker *stealingWorker) run(id int, deqs *([]*concurrent.Dequeue), doneTaskCount *int, totalTask int, wg *sync.WaitGroup) {
	task := (*deqs)[id].PopBottom()
	for {
		for task != nil {
			func_objs.ExecuteRunnable((*task).Func, (*task).Arg, doneTaskCount)
			task = (*deqs)[id].PopBottom()
		} 
		for task == nil {
			if *doneTaskCount != totalTask {
				victim := rand.Intn(len(*deqs))

				if !(*deqs)[victim].IsEmpty() {
					task = (*deqs)[victim].PopTop()
				}
			} else {
				// if there is no more works in global, stealing worker should stop stealing from others
				wg.Done()
				return
			}
		}
	}
	wg.Done()
}


func RunStealing(config Config) {
	scheduler := NewStealingScheduler(config.ThreadCount, config)
	// assign image tasks to each workers' dequeue
	for i := 0; i < len(*scheduler.Tasks); i++ {
		imageTask := (*scheduler.Tasks)[i]
		var ctx *func_objs.WorkerContext
		ctx = &func_objs.WorkerContext{Mutex:&sync.Mutex{}, Group:&sync.WaitGroup{}, Effects:imageTask.GetEffects()}
		scheduler.PushTask(i % config.ThreadCount, imageTask.ExecuteEffectsSteal, ctx)
	}
	
	scheduler.Run()
}