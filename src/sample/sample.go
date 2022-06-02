package main

import (
	"os"
	"proj2/png"
)

func main() {

	/******
		The following code shows you how to work with PNG files in Golang.
	******/

	//Assumes the user specifies a file as the first argument
	filePath := os.Args[1]

	//Loads the png image and returns the image or an error
	pngImg, err := png.Load(filePath)

	if err != nil {
		panic(err)
	}

	//Performs a grayscale filtering effect on the image
	pngImg.Grayscale()
	

	//Saves the image to a new file
	err = newImg.Save("test_gray2.png")

	//Checks to see if there were any errors when saving.
	if err != nil {
		panic(err)
	}

}
