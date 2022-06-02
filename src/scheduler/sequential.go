package scheduler

import (
	"strings"
	"proj2/iohandler"
)

func RunSequential(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	infoMaps := iohandler.ReadTxt(dataDirs)
	imageTasks := iohandler.CreateImageTasks(infoMaps, dataDirs)
	
	for _, imageTask := range imageTasks {
		effects := imageTask.GetEffects()
		maxY := imageTask.GetMaxY()
		minY := imageTask.GetMinY()
		waitgroup := make(chan bool)

		for _, effect := range effects {
			go imageTask.ExecuteEffects(effect, minY, maxY, waitgroup)
			<- waitgroup
		}
		imageTask.Swap()
		imageTask.SaveToPNG(imageTask.GetDataDir(), imageTask.GetOutPath())
	}
}
