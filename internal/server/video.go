package server

import (
	"context"
	"io"

	"github.com/google/uuid"
	"github.com/nei7/ntube/internal/db"
	"github.com/nei7/ntube/internal/service"
	"github.com/nei7/ntube/pkg"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type VideoServer struct {
	pkg.UnimplementedVideoUploadServiceServer
	uploadService service.VideoUpload
	tokenManager  service.TokenManager
	logger        *zap.Logger
}

func NewVideoServer(
	uploadService service.VideoUpload,
	tokenManager service.TokenManager,
	logger *zap.Logger,
) *VideoServer {
	return &VideoServer{
		pkg.UnimplementedVideoUploadServiceServer{},
		uploadService,
		tokenManager,
		logger,
	}
}

func (s *VideoServer) UploadVideo(stream pkg.VideoUploadService_UploadVideoServer) error {
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

	ownerId, err := uuid.Parse(userid)
	if err != nil {
		return status.Error(codes.Unauthenticated, "invalid user id")
	}

	file, id, err := s.uploadService.NewUpload()
	if err != nil {
		return status.Error(codes.Internal, "cannot save image to the store")
	}

	defer file.Close()

	req, err := stream.Recv()
	if err != nil {
		return err
	}

	videoSize, err := readChunks(stream, file)

	err = s.uploadService.Process(stream.Context(), db.CreateVideoParams{
		Title:       req.GetInfo().Title,
		Description: req.GetInfo().Description,
		Path:        id,
		Thumbnail:   id,
		OwnerID:     ownerId,
	})
	if err != nil {

		return err
	}

	err = stream.SendAndClose(&pkg.UploadVideoResponse{
		Id:   id,
		Size: videoSize,
	})

	if err != nil {
		return status.Errorf(codes.Unknown, "cannot send response: %v", err)
	}

	return nil
}

// 5GB
const MAX_VIDEO_SIZE = 5 * 8 * 1024 * 1024 * 1024

// 50 MB
const MAX_CHUNK_SIZE = 50 * 8 * 1024 * 1024

func readChunks(stream pkg.VideoUploadService_UploadVideoServer, src io.Writer) (uint64, error) {
	size := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}

		switch stream.Context().Err() {
		case context.Canceled:
			return 0, status.Error(codes.Aborted, "uploading canceled")
		case context.DeadlineExceeded:
			return 0, status.Error(codes.DeadlineExceeded, "deadline exceeded")
		}

		chunk := req.GetChunkData()

		n, err := src.Write(chunk)
		if err != nil {
			return 0, status.Error(codes.Internal, "failed to upload")
		}

		size += n

		if size > MAX_VIDEO_SIZE {
			return 0, status.Error(codes.InvalidArgument, "video is too large")
		}

	}

	return uint64(size), nil
}
