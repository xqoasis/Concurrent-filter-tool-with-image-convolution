package main

import (
	"proj3/png"
)

func main() {

	/******
		The following code shows you how to work with PNG files in Golang.
	******/

	//Assumes the user specifies a file as the first argument
	filePath := "/Users/xqqqq/github-classroom/mpcs-jh/project-2-xqoasis/proj3/sample/test_img.png"

	//Loads the png image and returns the image or an error
	pngImg, err := png.Load(filePath, nil)

	if err != nil {
		panic(err)
	}

	//Performs a grayscale filtering effect on the image
	// pngImg.Grayscale()
	pngImg.OtherEffect("B", pngImg.Bounds.Min.Y, pngImg.Bounds.Max.Y)

	//Saves the image to a new file
	err = pngImg.Save("new_B.png")

	//Checks to see if there were any errors when saving.
	if err != nil {
		panic(err)
	}

}
