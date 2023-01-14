package service

import (
	"fmt"
	"math"
	"os/exec"
	"strconv"
	"strings"
)

type FfpmegService interface {
	ExtractHLS(input string, output string) error
	DoScreenshot(input string, output string) error
	GetVideoDuration(path string) (float64, error)
}

type ffmpegService struct {
}

func NewFfpmegService() *ffmpegService {
	return &ffmpegService{}
}

func (s *ffmpegService) ExtractHLS(input string, output string) error {
	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		"-i", input,
		// video size
		"-vf", "scale=-1:720",
		// video codec
		"-c:v", "libx264",
		// video quality 51 is the worst quality and 1 the best
		"-crf", "21",
		// audio codec
		"-c:a", "libmp3lame",
		// format
		"-f", "mp4",
		output,
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *ffmpegService) DoScreenshot(input string, output string) error {

	duration, err := s.GetVideoDuration(input)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		// input
		"-i", input,
		// seek the position to the specified timestamp
		"-ss", fmt.Sprintf("00:%d:%d", int(math.Floor((duration/3)/60)), int(math.Floor(duration/3))),
		// only one frame
		"-vframes", "1",
		// control output quality
		"-q:v", "2",
		// output
		output,
	)

	_, err = cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (s *ffmpegService) GetVideoDuration(path string) (float64, error) {
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)

	o, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	duration, err := strconv.ParseFloat(strings.Replace(string(o), "\n", "", 1), 64)
	if err != nil {
		return 0, err
	}

	return duration, nil
}
