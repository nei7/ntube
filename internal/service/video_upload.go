package service

import (
	"context"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/nei7/ntube/internal/db"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var defaultFilePerm = os.FileMode(0664)

type VideoUpload interface {
	NewUpload() (*os.File, string, error)
	ThumbnailPath() string
	VideoPath() string
	Process(ctx context.Context, video db.CreateVideoParams) error
}

type videoUpload struct {
	thumbnailPath string
	videoPath     string
	ffmpegService FfpmegService
	videoService  VideoService
}

func NewVideoUpload(thumbnailPath, videoPath string, fmffmpegService FfpmegService, videoService VideoService) *videoUpload {
	return &videoUpload{thumbnailPath, videoPath, fmffmpegService, videoService}
}

func (s *videoUpload) NewUpload() (*os.File, string, error) {
	id := uuid.New()

	file, err := os.OpenFile(s.binPath(id.String()), os.O_CREATE|os.O_WRONLY, defaultFilePerm)

	return file, id.String(), err
}

func (s *videoUpload) Process(ctx context.Context, video db.CreateVideoParams) error {
	vid, err := s.videoService.Create(ctx, video)
	if err != nil {
		return err
	}

	binPath := s.binPath(vid.Path)

	err = s.ffmpegService.ExtractHLS(binPath, path.Join(s.videoPath, video.Path+".mp4"))
	if err != nil {
		return status.Error(codes.Internal, "failed to convert video")
	}

	err = s.ffmpegService.DoScreenshot(binPath, path.Join(s.thumbnailPath, video.Path+".jpg"))
	if err != nil {
		return err
	}

	return nil
}

func (s *videoUpload) binPath(id string) string {
	return path.Join(os.TempDir(), id)
}

func (s *videoUpload) ThumbnailPath() string {
	return s.thumbnailPath
}

func (s *videoUpload) VideoPath() string {
	return s.videoPath
}
