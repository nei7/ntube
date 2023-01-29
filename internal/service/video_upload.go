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

var defaultFilePerm = os.FileMode(0777)

type VideoUpload interface {
	NewUpload() (*os.File, string, error)
	Process(ctx context.Context, video db.CreateVideoParams) error
}

type videoUpload struct {
	uploadPath    string
	ffmpegService FfpmegService
	videoService  VideoService
}

func NewVideoUpload(uploadPath string, fmffmpegService FfpmegService, videoService VideoService) *videoUpload {
	return &videoUpload{uploadPath, fmffmpegService, videoService}
}

func (s *videoUpload) NewUpload() (file *os.File, id string, err error) {
	id = uuid.New().String()
	err = os.Mkdir(path.Join(s.uploadPath, id), defaultFilePerm)
	if err != nil {
		return
	}

	file, err = os.OpenFile(path.Join(s.uploadPath, id, "upload"), os.O_CREATE|os.O_WRONLY, defaultFilePerm)
	return
}

func (s *videoUpload) Process(ctx context.Context, video db.CreateVideoParams) error {
	defer otelSpan(ctx, "VideoUpload.Process").End()

	vid, err := s.videoService.Create(ctx, video)
	if err != nil {
		return err
	}

	uploadDir := path.Join(s.uploadPath, vid.Path)
	input := path.Join(uploadDir, "upload")

	err = s.ffmpegService.ProcessVideo(input, uploadDir)

	if os.Remove(input) != nil {
		return status.Error(codes.Internal, "failed to remove video")
	}

	if err != nil {
		return status.Error(codes.Internal, "failed to convert video")
	}

	return nil
}
