package ffmpeg

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

func ExtractHLS(file_path, filename string) error {
	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		"-i", path.Join(file_path, "tmp", filename+".tmp"),
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

		path.Join(file_path, "mp4", filename+".mp4"),
	)

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func DoScreenshot(file_path, filename string) error {
	tmpPath := path.Join(file_path, "tmp", filename+".tmp")

	duration, err := getVideoDuration(tmpPath)
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
		path.Join(file_path, "thumbnail", filename+".jpg"),
	)

	fmt.Println(cmd.String())
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func getVideoDuration(path string) (float64, error) {
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
