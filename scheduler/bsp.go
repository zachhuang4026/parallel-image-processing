package scheduler

import (
	"sync"
	"strings"
	"proj2/iohandler"
	"proj2/png"
)

type bspWorkerContext struct {
	// Define the necessary fields for your implementation
	Mutex             *sync.Mutex
	Cond              *sync.Cond
	WgContext         *sync.WaitGroup
	Tasks             *([]png.ImageTask)
	CurrentTaskIdx    int
	CurrentEffectIdx  int
	ThreadCount       int
	DoneFlag          *[]bool
	DoneCount         int
	EffectFlag        *[]bool
	EffectCount       int
}

func NewBSPContext(config Config) *bspWorkerContext {
	//Initialize the context
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	var wg sync.WaitGroup
	dataDirs := strings.Split(config.DataDirs, "+")
	infoMaps := iohandler.ReadTxt(dataDirs)
	tasks := iohandler.CreateImageTasks(infoMaps, dataDirs)
	doneFlag := make([]bool, config.ThreadCount)
	effectFlag := make([]bool, 4)
	return &bspWorkerContext{Mutex: &mu, Cond: cond, WgContext: &wg, Tasks: &tasks,ThreadCount: config.ThreadCount, DoneFlag: &doneFlag, EffectFlag: &effectFlag}
}

func RunBSPWorker(id int, ctx *bspWorkerContext) {
	// Implement the BSP model here.
	// No additional loops can be used in this implementation. T
	// This goes to calling other functions. No other called
	// function you define or are using can have looping being done for you.
	
	for {
		// if all image task be done, execute the program
		if ctx.CurrentTaskIdx == len(*ctx.Tasks) {
			break
		}
		// idx = 0, do #0 region; idx = 1, do #1 region
		imgTask := (*ctx.Tasks)[ctx.CurrentTaskIdx]
		turns := imgTask.GetMaxY() / ctx.ThreadCount
		effects := imgTask.GetEffects()	
		
		imgTask.ExecuteEffectsBSP(effects[ctx.CurrentEffectIdx], turns*(id), turns*(id+1))
		
		ctx.Mutex.Lock()
		(*ctx.DoneFlag)[id] = true
		ctx.DoneCount++
		ctx.Mutex.Unlock()
		
		// all workers done one effect
		if ctx.DoneCount == ctx.ThreadCount {
			ctx.Mutex.Lock()
			ctx.EffectCount++
			ctx.Mutex.Unlock()

			// still have other effects need to be done
			if ctx.EffectCount != len(effects) {
				ctx.Mutex.Lock()
				(*ctx.EffectFlag)[ctx.CurrentEffectIdx] = true
				ctx.CurrentEffectIdx++
				(*ctx.DoneFlag) = make([]bool, ctx.ThreadCount)
				ctx.DoneCount = 0
				ctx.Cond.Broadcast()
				ctx.Mutex.Unlock()
				
			} else {
				// last global sync: store the image task to png
				ctx.Mutex.Lock()
				imgTask.Swap()
				imgTask.SaveToPNG(imgTask.GetDataDir(), imgTask.GetOutPath())
	
				// reset DoneCount, Doneflag, CurrentEffectIdx, EffectCount,
				ctx.DoneCount = 0
				ctx.CurrentEffectIdx = 0
				ctx.EffectCount = 0
				(*ctx.DoneFlag) = make([]bool, ctx.ThreadCount)
				(*ctx.EffectFlag) = make([]bool, 4)
	
				// move to next imageTask and wake all goroutines up
				ctx.CurrentTaskIdx++
				ctx.Cond.Broadcast()
				ctx.Mutex.Unlock()
			}
		} else {
			ctx.Mutex.Lock()
			for (*ctx.DoneFlag)[id] == true {
				ctx.Cond.Wait()
			}
			ctx.Mutex.Unlock()
		}
	}
}
