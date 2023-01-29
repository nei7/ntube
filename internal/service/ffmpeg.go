package service

import (
	"fmt"
	"math"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

type FfpmegService interface {
	ProcessVideo(input, output_dir string) error
}

type ffmpegService struct {
}

func NewFfpmegService() *ffmpegService {
	return &ffmpegService{}
}

func (s *ffmpegService) ProcessVideo(input, output_dir string) (err error) {
	err = s.createThumbnail(input, path.Join(output_dir, "thumbnail1.jpg"), 3)
	if err != nil {
		return
	}

	return s.ConvertVideo(input, output_dir)
}

func (s *ffmpegService) ConvertVideo(input, output string) error {
	args := ffmpeg.KwArgs{
		"c:v": "libx264",
		"crf": "21",
		"c:a": "libmp3lame",
		"f":   "mp4",
	}

	qualities := map[string]string{
		"720p": "scale=-1:720",
		"360p": "scale=-1:360",
	}

	videoInput := ffmpeg.Input(input)

	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	for name, v := range qualities {
		go func(name, v string) {
			args["vf"] = v
			err := videoInput.Output(path.Join(output, name+".mp4"), args).Run()
			if err != nil {
				errChan <- err
			}
			wg.Done()
		}(name, v)
	}

	wg.Add(2)

	wg.Wait()

	close(errChan)

	return <-errChan
}

func (s *ffmpegService) createThumbnail(input string, output string, position float64) error {
	duration := s.getVideoDuration(input)
	position = duration / position
	return ffmpeg.Input(input).Output(output, ffmpeg.KwArgs{
		"ss":      fmt.Sprintf("00:%d:%d", int(math.Floor(position/60)), int(math.Floor(position))),
		"vframes": "1",
		"q:v":     "2",
	}).OverWriteOutput().Run()
}

func (s *ffmpegService) getVideoDuration(path string) float64 {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	output, err := cmd.Output()
	if err != nil {
		return 0
	}
	duration, err := strconv.ParseFloat(strings.Replace(string(output), "\n", "", 1), 64)
	if err != nil {
		return 0
	}
	return duration
}
