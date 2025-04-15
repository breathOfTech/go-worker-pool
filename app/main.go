package main

import (
	"fmt"
	"streamer"
)

func main() {
	// Define number of workers and jobs

	const (
		numbJobs   = 1
		numWorkers = 2
	)

	// create channels for work and results
	notifyChannel := make(chan streamer.ProcessingMessage, numbJobs)
	defer close(notifyChannel)

	videoQueue := make(chan streamer.VideoProcessingJob, numbJobs)
	defer close(videoQueue)

	// Get a worker pool.
	wp := streamer.New(videoQueue, numWorkers)
	fmt.Println("wp:", wp)

	// start the worker pool.
	wp.Run()
	fmt.Println("worker pool started, press enter to continue..")
	fmt.Scanln()

	// create 4 videos to send to the worker pool.
	video := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChannel, nil)

	// send the videos to the worker pool.
	videoQueue <- streamer.VideoProcessingJob{Video: *video}

	// Print out results.
	for i := 1; i <= numbJobs; i++ {
		msg := <-notifyChannel
		fmt.Println("i:", i, "msg: ", msg)
	}

	fmt.Println("Done!")

}
