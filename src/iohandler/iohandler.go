package iohandler

import (
    "encoding/json"
	"image"
    "fmt"
	"bufio"
	"os"
	"log"
	png "proj2/png"
)

func swap(px, py *image.RGBA64) {
	tempx := *px
	tempy := *py
	*px = tempy
	*py = tempx
}

func ReadTxt(dataDirs []string) []map[string]interface{} {
	file, err := os.Open("../data/effects.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// read effects.txt line by line
	scanner := bufio.NewScanner(file)
	infoMaps := []map[string]interface{}{}
	for scanner.Scan() { 
		// convert string to map
		// ref: https://stackoverflow.com/questions/51189959/convert-output-string-to-map-with-golang
        newLine := scanner.Text()
		infoMap := map[string]interface{}{}
    	if err := json.Unmarshal([]byte(newLine), &infoMap); err != nil {
        	panic(err)
    	}
		infoMaps = append(infoMaps, infoMap)
	}
	if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
	return infoMaps
}

func CreateImageTasks(infoMaps []map[string]interface{}, dataDirs []string) []png.ImageTask {
	var imageTasks []png.ImageTask

	for _, infoMap := range infoMaps {
		inPath := infoMap["inPath"].(string)
    	outPath := infoMap["outPath"].(string)
		effects := infoMap["effects"].([]interface{})

		for _, dataDir := range dataDirs {
			filePath := fmt.Sprintf("../data/in/%v/%v", dataDir, inPath)
			pngImg, err := png.Load(filePath)
	
			if err != nil {
				panic(err)
			}
			pngImg.SetDataDir(dataDir)
			pngImg.SetInPath(string(inPath))
			pngImg.SetOutPath(string(outPath))
			pngImg.SetEffects(effects)
			
			imageTasks = append(imageTasks, *pngImg)
		}
	}
	return imageTasks
}