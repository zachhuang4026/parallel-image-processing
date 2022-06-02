package scheduler

import (
	"sync"
	"strings"
	"math/rand"
	"proj3/workqueue"
	"proj3/png"
	"proj3/iohandler"
	"proj3/func_objs"
)

type BalancingScheduler struct {
	Mutex *sync.Mutex
	Group *sync.WaitGroup
	Queues *([]*workqueue.Queue)
	Random int
	Tasks *([]png.ImageTask)
	ThreadCount int
	Workers *([]balancingWorker)
	DoneTaskCount *int
}

func NewBalancingScheduler(capacity int, config Config) *BalancingScheduler {
	queues := make([](*workqueue.Queue), capacity)

	for i := 0; i < config.ThreadCount; i++ {
		queue := workqueue.NewQueue()
		queues[i] = queue
	} 
	
	dataDirs := strings.Split(config.DataDirs, "+")
	infoMaps := iohandler.ReadTxt(dataDirs)
	tasks := iohandler.CreateImageTasks(infoMaps, dataDirs)
	workers := make([]balancingWorker, capacity)
	doneTaskCount := 0

	return &BalancingScheduler{Mutex: &sync.Mutex{}, Group: &sync.WaitGroup{}, ThreadCount: config.ThreadCount,
		Queues: &queues, Tasks: &tasks, Workers: &workers, DoneTaskCount: &doneTaskCount}
}

func (scheduler *BalancingScheduler) PushTask(idx int, task func_objs.Runnable, taskArg interface{}) {
	(*scheduler.Queues)[idx].Enq(task, taskArg)
}

func (scheduler *BalancingScheduler) Run(){
	scheduler.Group.Add(scheduler.ThreadCount)
	for i := 0; i < scheduler.ThreadCount; i++ {
		go (*scheduler.Workers)[i].run(i, scheduler.Queues, scheduler.DoneTaskCount, len(*scheduler.Tasks), scheduler.Group, scheduler.Mutex)
	}
	scheduler.Wait()
}

func (scheduler *BalancingScheduler) Wait() {
	scheduler.Group.Wait()
}

type balancingWorker struct {
	queue    *workqueue.Queue
}

func NewBalancingWorker() *balancingWorker {
	queue := workqueue.NewQueue()
	return &balancingWorker{queue: queue}
}

func (worker *balancingWorker) run(id int, queues *([]*workqueue.Queue), doneTaskCount *int, totalTask int, wg *sync.WaitGroup, mutex *sync.Mutex) {	
	//*doneTaskCount != totalTask 
	for {
		task := (*queues)[id].Deq()
		if task != nil {
			func_objs.ExecuteRunnable((*task).Func, (*task).Arg, doneTaskCount)
		} 
			
		// after dequeing, try to balance work load from others
		size := (*queues)[id].Size()

		if rand.Intn(size+1) == size {
			victim := rand.Intn(len(*queues))
			var min int
			var max int
			if victim <= id {
				min = victim
				max = id
			} else {
				min = id
				max = victim
			}
			//mutex.Lock()
			(*queues)[min].Lock()
			(*queues)[max].Lock()
			balance((*queues)[min], (*queues)[max])
			(*queues)[max].UnLock()
			(*queues)[min].UnLock()
			//mutex.Unlock()
			
		} 
		
	} 
	wg.Done()
}

func balance(q0 *workqueue.Queue, q1 *workqueue.Queue) {
	var qMin workqueue.Queue
	var qMax workqueue.Queue
	if (*q0).Size() < (*q1).Size() {
		qMin = *q0
	} else {
		qMin = *q1
	}

	if (*q0).Size() < (*q1).Size() {
		qMax = *q1
	} else {
		qMax = *q0
	}

	diff := qMax.Size() - qMin.Size()

	if diff > 1 {
		for qMax.Size() > qMin.Size() {
			task := qMax.Deq()
			qMin.Enq((*task).Func, (*task).Arg)
		}
	}
}


func RunBalance(config Config) {
	scheduler := NewBalancingScheduler(config.ThreadCount, config)
	
	for i := 0; i < len(*scheduler.Tasks); i++ {
		imageTask := (*scheduler.Tasks)[i]
		var ctx *func_objs.WorkerContext
		ctx = &func_objs.WorkerContext{Mutex:&sync.Mutex{}, Group:&sync.WaitGroup{}, Effects:imageTask.GetEffects()}
		scheduler.PushTask(i % config.ThreadCount, imageTask.ExecuteEffectsSteal, ctx)
	}
	scheduler.Run()
}