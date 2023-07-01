package scheduler

import (
	"fmt"
	"proj3/png"
	"proj3/utils"
	"strings"
)



func ImgTaskHandler(imgTask *png.ImageTask) {
	// Assumes the user specifies a file as the first argument
	bounds := imgTask.Bounds
	for _, effect := range imgTask.CntTask.Effects {
		if effect == "G" {
			imgTask.Grayscale(bounds.Min.Y, bounds.Max.Y)
		}else {
			imgTask.OtherEffect(effect, bounds.Min.Y, bounds.Max.Y)
		}
		imgTask.SwapPointer()
	}
	imgTask.SwapPointer()
	outDirPath := fmt.Sprint(utils.DATADIR + "/out/", imgTask.CntTask.Dir, "/", imgTask.CntTask.OutPath)
	err := imgTask.Save(outDirPath)
	if err != nil {
		panic(err)
	}
}
func RunSequential(config Config) {
	dataDirs := strings.Split(config.DataDirs, "+")
	for _, dir := range dataDirs {
		decoder := utils.Reader()
		for {
			imgTask := utils.ImageTaskGeneretor(decoder, dir)
			if imgTask == nil {
				return
			}
			ImgTaskHandler(imgTask)
		}
	}
}