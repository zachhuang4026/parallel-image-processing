package main

import (
	"fmt"
	"os"
	"proj2/scheduler"
	"strconv"
	"time"
)

const usage = "Usage: editor data_dir mode [number of threads]\n" +
	"data_dir = The data directory to use to load the images.\n" +
	"mode     = (bsp) run the BSP mode, (pipeline) run the pipeline mode\n" +
	"[number of threads] = Runs the parallel version of the program with the specified number of threads.\n"

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		return
	}
	config := scheduler.Config{DataDirs: "", Mode: "", ThreadCount: 0}
	config.DataDirs = os.Args[1]

	if len(os.Args) >= 3 {
		config.Mode = os.Args[2]
		threads, _ := strconv.Atoi(os.Args[3])
		config.ThreadCount = threads
	} else {
		config.Mode = "s"
	}
	start := time.Now()
	scheduler.Schedule(config)
	end := time.Since(start).Seconds()
	fmt.Printf("%.2f\n", end)

}
