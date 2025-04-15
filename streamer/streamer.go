package streamer

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type ProcessingMessage struct {
	ID         int
	Successful bool
	Message    string
	OutputFile string
}

type VideoProcessingJob struct {
	Video Video
}

type Processor struct {
	Engine Encoder
}

type Video struct {
	ID           int
	InputFile    string
	OutputDir    string
	EncodingType string
	NotifyChan   chan ProcessingMessage
	Encoder      Processor
	Options      *VideoOptions
}

type VideoOptions struct {
	RenameOutput    bool
	SegmentDuration int
	MaxRate1080p    string
	MaxRate720p     string
	MaxRate480p     string
}

func (vd *VideoDispatcher) NewVideo(id int, input, output, encodingType string, notifyChan chan ProcessingMessage, ops *VideoOptions) *Video {
	if ops == nil {
		ops = &VideoOptions{}
	}
	fmt.Println("NewVideo: New Video created: ", id, input)
	return &Video{
		ID:           id,
		InputFile:    input,
		OutputDir:    output,
		EncodingType: encodingType,
		NotifyChan:   notifyChan,
		Encoder:      vd.Processsor,
		Options:      ops,
	}
}

func (v *Video) encode() {
	var fileName string

	switch v.EncodingType {
	case "mp4":
		{
			// encode the video
			fmt.Println("v.encode(): About to encode to MP4", v.ID)
			name, err := v.encodeToMP4()

			if err != nil {
				// send inforamtion to the notifyChan
				v.sendToNotifyChan(false, "", fmt.Sprintf("encode failed for %d: %s", v.ID, err.Error()))
				return
			}

			fileName = fmt.Sprintf("%s.mp4", name)
		}
	default:
		{
			fmt.Println("v.encode(): error trying to encode video", v.ID)
			v.sendToNotifyChan(false, "", fmt.Sprintf("error processing for %d: invalid enoding type", v.ID))
			return
		}
	}

	fmt.Println("v.encode(): sendig success messsaged for video if", v.ID, "to notifyChan")
	v.sendToNotifyChan(true, fileName, fmt.Sprintf("video id %d processed and saved as %s", v.ID, fmt.Sprintf("%s/%s", v.OutputDir, fileName)))
}

func (v *Video) encodeToMP4() (string, error) {
	baseFileName := ""

	fmt.Println("v.encodeToMP4: about to try to encode video id", v.ID)
	if !v.Options.RenameOutput {
		// Get the base file name
		b := path.Base(v.InputFile)
		baseFileName = strings.TrimSuffix(b, filepath.Ext(b))
	} else {
		// TODO: Generate random file name
	}
	err := v.Encoder.Engine.EncodeToMP4(v, baseFileName)
	if err != nil {
		return "", err
	}

	fmt.Println("v.encodeToMP4: successfully encoded video id", v.ID)

	return baseFileName, nil
}
func (v *Video) sendToNotifyChan(successful bool, fileName, message string) {
	fmt.Println("v.sendToNotifyChan: sending message to notifyChan for video id", v.ID)
	v.NotifyChan <- ProcessingMessage{
		ID:         v.ID,
		Successful: successful,
		Message:    message,
		OutputFile: fileName,
	}
}

func New(jobQueue chan VideoProcessingJob, maxWorkers int) *VideoDispatcher {
	fmt.Println("New: creating worker pool")
	workerPool := make(chan chan VideoProcessingJob, maxWorkers)

	// TODO: implement processor logic
	var e VideoEncoder
	p := Processor{
		Engine: &e,
	}

	return &VideoDispatcher{
		WorkerPool: workerPool,
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		Processsor: p,
	}
}
