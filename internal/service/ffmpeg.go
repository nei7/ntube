package service

import (
	"fmt"
	"math"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type FfpmegService interface {
	ExtractHLS(filename string) error
	DoScreenshot(filename string) error
	GetVideoDuration(path string) (float64, error)
}

type ffmpegService struct {
	basePath string
}

func NewFfpmegService(basePath string) *ffmpegService {
	return &ffmpegService{basePath}
}

func (s *ffmpegService) ExtractHLS(filename string) error {
	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		"-i", path.Join(s.basePath, "tmp", filename+".tmp"),
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
		path.Join(s.basePath, "mp4", filename+".mp4"),
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func (s *ffmpegService) DoScreenshot(filename string) error {
	tmpPath := path.Join(s.basePath, "tmp", filename+".tmp")

	duration, err := s.GetVideoDuration(tmpPath)
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		// input
		"-i", tmpPath,
		// seek the position to the specified timestamp
		"-ss", fmt.Sprintf("00:%d:%d", int(math.Floor((duration/3)/60)), int(math.Floor(duration/3))),
		// only one frame
		"-vframes", "1",
		// control output quality
		"-q:v", "2",
		// output
		path.Join(s.basePath, "thumbnail", filename+".jpg"),
	)

	err = cmd.Start()
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
