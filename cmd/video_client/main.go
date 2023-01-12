package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	"github.com/nei7/gls/internal/service"
	"github.com/nei7/gls/pkg/video"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {

	var videoPath string

	flag.StringVar(&videoPath, "v", "", "video path")

	flag.Parse()

	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("failed to load .env file", err)
	}

	conn, err := grpc.Dial(":3001", grpc.WithInsecure())
	if err != nil {
		log.Fatal("error while connecting", err)
	}
	videoClient := video.NewVideoUploadServiceClient(conn)

	f, err := os.Open(videoPath)
	if err != nil {
		log.Fatal("error while opening file", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tokenManager := service.NewTokenManager(viper.GetString("JWT_KEY"))

	jwt, err := tokenManager.NewJWT("7b129eba-8b82-4ff0-9751-a78fb1868993", time.Now().Add(time.Hour*2).Unix())
	if err != nil {
		log.Fatal("failed to generate jwt token")
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "authorization", jwt)

	stream, err := videoClient.UploadVideo(ctx)
	if err != nil {
		log.Fatal("failed to establish connection with video server", err)
	}

	err = stream.Send(&video.UploadVideoRequest{
		Data: &video.UploadVideoRequest_Info{
			Info: &video.VideoInfo{
				Title: "video",
			},
		},
	})
	if err != nil {
		log.Fatalf("can't send image to server")
	}

	reader := bufio.NewReader(f)
	buf := make([]byte, 1024)

	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("can't send chunk data to server", err)
		}

		err = stream.Send(&video.UploadVideoRequest{
			Data: &video.UploadVideoRequest_ChunkData{
				ChunkData: buf[:n],
			},
		})
		if err != nil {
			log.Fatal("can't send chunk data to server", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("can't get response", err)
	}

	log.Printf("Image uploaded successfully - ID: %s, size: %d \n", res.GetId(), res.GetSize())
}
