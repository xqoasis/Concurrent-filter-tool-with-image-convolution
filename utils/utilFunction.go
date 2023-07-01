package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"proj3/png"
)
var REMOTE_DATADIR = "/dev/shm/data"
var LOCAL_DATADIR = "../data"
var DATADIR = LOCAL_DATADIR

func Reader() *json.Decoder{
	effectsPathFile := fmt.Sprintf(DATADIR + "/effects.txt")
	effectsFile, err := os.Open(effectsPathFile)
	if err != nil {
		log.Fatal(err)
	}
	return json.NewDecoder(effectsFile)
}

func TaskGeneretor(decoder *json.Decoder, dir string) *png.Task{
	// generate each task
	var task png.Task
	if err := decoder.Decode(&task); err == io.EOF {
		// println("End of effects.txt")
		return nil
	}else if err != nil {
		fmt.Println(err)
		return nil
	}
	task.Dir = dir
	return &task
}

func ImageTaskGeneretor(decoder *json.Decoder, dir string) *png.ImageTask{
	// generate each ImageTask
	task := TaskGeneretor(decoder, dir)
	if task == nil {
		return nil
	}
	inDirPath := fmt.Sprint(DATADIR + "/in/", task.Dir, "/", task.InPath)
	imagetask, err := png.Load(inDirPath, task)
	if err != nil {
		panic(err)
	}
	return imagetask
}