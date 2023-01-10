package grpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/nei7/gls/pkg/video"
)

type VideoServer struct {
	video.UnimplementedVideoUploadServiceServer
	storePath string
}

func (s *VideoServer) UploadVideo(stream video.VideoUploadService_UploadVideoServer) error {
	req, err := stream.Recv()
	if err != nil {
		return err
	}
	title := req.GetInfo().GetTitle()

	videoBuf := bytes.Buffer{}
	videoSize := 0

	for {

		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		switch stream.Context().Err() {
		case context.Canceled:
			return errors.New("request canceled")
		case context.DeadlineExceeded:
			return errors.New("deadline exceeded")
		default:
			return nil
		}

		chunk := req.GetChunkData()
		size := len(chunk)

		videoSize += size

		if _, err = videoBuf.Write(chunk); err != nil {
			return err
		}

	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	file, err := os.Create(path.Join(s.storePath, fmt.Sprintf("%s.%s", id.String(), "tmp")))
	if err != nil {
		return nil
	}

	_, err = videoBuf.WriteTo(file)
	if err != nil {
		return err
	}

	return nil
}
