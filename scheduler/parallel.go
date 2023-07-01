package scheduler

import (
	"proj3/concurrent"
	"proj3/png"
	"proj3/utils"
	"strings"
)

// Making the imgTaskWrap runnable, should define a new struct in this file
type imgTaskWrap struct {
	imgTask *png.ImageTask
}
// Making the imgTaskWrap runnable
func (imgTaskWrap *imgTaskWrap) Run() {
	ImgTaskHandler(imgTaskWrap.imgTask)
}

func RunParallel(config Config) {
	var futures []concurrent.Future
	var exeService concurrent.ExecutorService
	if config.Mode == "steal" {
		exeService = concurrent.NewWorkStealingExecutor(config.ThreadCount, 10)
	}else if config.Mode == "balance" {
		exeService = concurrent.NewWorkBalancingExecutor(config.ThreadCount, 10, 2)
	}

	dataDirs := strings.Split(config.DataDirs, "+")
	decoder := utils.Reader()
	for _, dir := range dataDirs {
		for {
			imgTask := utils.ImageTaskGeneretor(decoder, dir)
			if imgTask == nil {
				exeService.Shutdown()
				break
			} else {
				futures = append(futures, exeService.Submit(&imgTaskWrap{imgTask: imgTask}))
			}
			
		}
	}
	for _, future := range futures {
		future.Get()
	}
	return
}