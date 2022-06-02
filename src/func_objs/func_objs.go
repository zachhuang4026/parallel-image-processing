package func_objs

import (
	"sync"
)

type Runnable func(arg interface{})

type WorkerContext struct {
	Mutex *sync.Mutex
	Group *sync.WaitGroup
	Effects []string
}


func ExecuteRunnable(runnable Runnable, arg interface{}, doneTaskCount *int){
	ctx := arg.(*WorkerContext)
	runnable(ctx)
	*doneTaskCount++
}