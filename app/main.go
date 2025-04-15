package main

import (
	"fmt"
	"streamer"
)

func main() {
	// Define number of workers and jobs

	const (
		numbJobs   = 4
		numWorkers = 4
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

	// create 4 videos to send to the worker pool.

	// create a video that converts mp4 to web ready format
	video1 := wp.NewVideo(1, "./input/puppy1.mp4", "./output", "mp4", notifyChannel, nil)

	// create second video that should fail
	video2 := wp.NewVideo(2, "./input/bad.txt", "./output", "mp4", notifyChannel, nil)

	// create third video that converts mp4 to hls
	ops := &streamer.VideoOptions{
		RenameOutput:    true,
		SegmentDuration: 10,
		MaxRate1080p:    "1200k",
		MaxRate720p:     "600k",
		MaxRate480p:     "400k",
	}
	video3 := wp.NewVideo(3, "./input/puppy2.mp4", "./output", "hls", notifyChannel, ops)

	video4 := wp.NewVideo(4, "./input/puppy2.mp4", "./output", "mp4", notifyChannel, nil)

	// send the videos to the worker pool.
	videoQueue <- streamer.VideoProcessingJob{Video: *video1}
	videoQueue <- streamer.VideoProcessingJob{Video: *video2}
	videoQueue <- streamer.VideoProcessingJob{Video: *video3}
	videoQueue <- streamer.VideoProcessingJob{Video: *video4}

	// Print out results.
	for i := 1; i <= numbJobs; i++ {
		msg := <-notifyChannel
		fmt.Println("i:", i, "msg: ", msg)
	}

	fmt.Println("Done!")
}
