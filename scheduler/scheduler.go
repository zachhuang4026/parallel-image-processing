package scheduler

type Config struct {
	DataDirs string //Represents the data directories to use to load the images.
	Mode     string // Represents which scheduler scheme to use
	// If Mode == "s" run the sequential version
	// If Mode == "pipeline" run the pipeline version
	// If Mode == "bsp" run the pipeline version
	// These are the only values for Version
	ThreadCount int // Runs the parallel version of the program with the
	// specified number of threads (i.e., goroutines)
}

//Run the correct version based on the Mode field of the configuration value
func Schedule(config Config) {
	if config.Mode == "s" {
		RunSequential(config)
	} else if config.Mode == "pipeline" {
		RunPipeline(config)
	} else if config.Mode == "bsp" {
		ctx := NewBSPContext(config)
		var idx int
		for idx = 0; idx < config.ThreadCount-1; idx++ {
			go RunBSPWorker(idx, ctx)
		}
		RunBSPWorker(idx, ctx)
	} else {
		panic("Invalid scheduling scheme given.")
	}
}
