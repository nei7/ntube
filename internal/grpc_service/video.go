package grpc_service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/nei7/gls/internal/db"
	"github.com/nei7/gls/internal/ffmpeg"
	"github.com/nei7/gls/internal/service"
	"github.com/nei7/gls/pkg/video"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// 5GB
const MAX_VIDEO_SIZE = 5 * 8 * 1024 * 1024 * 1024

// 50 MB
const MAX_CHUNK_SIZE = 50 * 8 * 1024 * 1024

type VideoServer struct {
	video.UnimplementedVideoUploadServiceServer
	storePath    string
	svc          service.VideoService
	tokenManager service.TokenManager
	logger       *zap.Logger
}

func NewVideoServer(storePath string, svc service.VideoService, tokenManager service.TokenManager, logger *zap.Logger) *VideoServer {
	return &VideoServer{video.UnimplementedVideoUploadServiceServer{}, storePath, svc, tokenManager, logger}
}

func (s *VideoServer) UploadVideo(stream video.VideoUploadService_UploadVideoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	userid, err := s.tokenManager.Parse(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "invalid access token")
	}

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
			return status.Error(codes.Aborted, "uploading canceled")
		case context.DeadlineExceeded:
			return status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}

		chunk := req.GetChunkData()
		size := len(chunk)
		if size > MAX_CHUNK_SIZE {
			return status.Error(codes.InvalidArgument, "too large chunk")
		}

		videoSize += size
		if videoSize > MAX_VIDEO_SIZE {
			return status.Error(codes.InvalidArgument, "video is too large")
		}

		if _, err = videoBuf.Write(chunk); err != nil {
			return err
		}

	}
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	ownerId, err := uuid.Parse(userid)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid user id")
	}

	_, err = s.svc.Create(stream.Context(), db.CreateVideoParams{
		Path:    id.String(),
		OwnerID: ownerId,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	tmpPath := path.Join(s.storePath, "tmp", fmt.Sprintf("%s.%s", id.String(), "tmp"))
	file, err := os.Create(tmpPath)
	if err != nil {
		return status.Error(codes.Internal, "cannot save image to the store")
	}

	_, err = videoBuf.WriteTo(file)
	if err != nil {
		return status.Error(codes.Internal, "cannot save image to the store")
	}

	err = ffmpeg.Extract_HLS(s.storePath, id.String())
	if err != nil {
		return status.Error(codes.Internal, "failed to convert video")
	}

	return nil
}
