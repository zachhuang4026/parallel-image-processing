package scheduler

import (
	"strings"
	"proj2/png"
	"proj2/iohandler"
)


func RunPipeline(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	infoMaps := iohandler.ReadTxt(dataDirs)
	imageTasks := iohandler.CreateImageTasks(infoMaps, dataDirs)

	/* Push all imageTasks onto imgTaskStream */
	repeatFn := func(
		done <-chan interface{},
		imgTasks []png.ImageTask,
	) <-chan png.ImageTask {
		imgTaskStream := make(chan png.ImageTask)
		go func() {
			defer close(imgTaskStream)
			for i := 0; i < len(imgTasks); i++ {
				select {
				case <-done: 
					return
				case imgTaskStream <- imgTasks[i]:
				}
			}
		}()
		return imgTaskStream
	}
	/* Execute Effects on each ImageTask */
	effectDoer := func(
		done <-chan interface{}, 
		imgTaskStream <-chan png.ImageTask,
	) <-chan png.ImageTask {
		imgResultStream := make(chan png.ImageTask)
		
		go func() {
			defer close(imgResultStream)
			for imgTask := range imgTaskStream {
				maxY := imgTask.GetMaxY()
				turns := maxY / config.ThreadCount
				effects := imgTask.GetEffects()
				
				for _, effect := range effects {	
					// every effect launch new mini-workers
					// wait until all mini-workers done for the next effect
					waitgroup := make(chan bool)
					if maxY % config.ThreadCount == 0 {
						for i := 0; i < config.ThreadCount; i++ {
							go imgTask.ExecuteEffects(effect, turns*(i), turns*(i+1), waitgroup)
						}
					} else {
						for i := 0; i < config.ThreadCount-1; i++ {
							go imgTask.ExecuteEffects(effect, turns*(i), turns*(i+1), waitgroup)
						}
						go imgTask.ExecuteEffects(effect, turns*(config.ThreadCount-1), maxY, waitgroup)
					}
					for i := 0; i < config.ThreadCount; i++ {
						<- waitgroup
					}
				}
				imgTask.Swap()
				select {
				case <-done:
					return
				case imgResultStream <- imgTask:
				}
				
			}
		}()
		
		return imgResultStream
	}
	/* Fan-In: launch workers to execute EffectsDoer() */
	fanIn := func(
		done <-chan interface{},
		ttlTasks int,
		channels ...<-chan png.ImageTask,
	) <-chan png.ImageTask { 
		// counter how many imageTask been done
		counter := 0

		// initial a buffer channel with all true bool
		waitgroup := make(chan bool, config.ThreadCount)
		for i := 0; i < config.ThreadCount; i++ {
			waitgroup <- true
		}

		multiplexedStream := make(chan png.ImageTask)
		multiplex := func(c <-chan png.ImageTask, waitgroup chan bool, counter *int) { 
			for i := range c {
				select {
				case <-done:
					return
				case multiplexedStream <- i:
				}
				*counter++
			}
			// remove one buffer space for the next goroutine
			<- waitgroup  
			
		}

		for _, c := range channels {
			go multiplex(c, waitgroup, &counter)
		}
		
		// Wait for all the reads to complete
		go func() {
			defer close(multiplexedStream)
			// wait until all imgTask done by all goroutines
			for counter != ttlTasks {}
		}()

		return multiplexedStream // return subImgTask
	}
	savePNG := func(
		done <-chan interface{},
		resultImgStream <-chan png.ImageTask,
	) <-chan png.ImageTask{
		finishStream := make(chan png.ImageTask)
		go func() {
			defer close(finishStream)
			
			for resultImg := range resultImgStream {
				resultImg.SaveToPNG(resultImg.GetDataDir(), resultImg.GetOutPath())
			}
			select {
			case <-done:
				return
			case finishStream <- (<-resultImgStream):
			}
		}()
		return finishStream
	}


	done := make(chan interface{})
	defer close(done)

	imgTaskStream := repeatFn(done, imageTasks)

	numWorkers := config.ThreadCount
	workers := make([]<-chan png.ImageTask, numWorkers)
	for i := 0; i < numWorkers; i++ {
		workers[i] = effectDoer(done, imgTaskStream)
	}
	// safe image to PNG
	for finishImg := range savePNG(done, fanIn(done, len(imageTasks), workers...)) {
		finishImg.GetOutPath()
	}
}
