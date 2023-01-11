package ffmpeg

import (
	"os/exec"
	"path"
)

func Extract_HLS(file_path, filename string) error {
	cmd := exec.Command(
		"/usr/bin/ffmpeg",
		"-i", path.Join(file_path, "tmp", filename+".tmp"),
		// video size
		"-vf", "scale=720x?",
		// video codec
		"-c:v", "libx264",
		// video quality 51 is the worst quality and 1 the best
		"-crf", "21",
		// audio codec
		"-c:a", "libmp3lame",
		// format
		"-f", "hls",
		// slices the video and audio into segments with a duration of 6 seconds
		"-hls_time", "6",
		path.Join(file_path, "hls", filename+".m3u8"),
	)

	cmd.Stdin = nil
	cmd.Stderr = nil

	err := cmd.Start()
	if err != nil {
		return err
	}

	return nil
}
