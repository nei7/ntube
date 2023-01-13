package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	esv7api "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/nei7/ntube/internal"
	"github.com/nei7/ntube/internal/datastruct"
)

type Video struct {
	client *esv7.Client
	index  string
}

type indexedVideo struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	Path        string `json:"path"`
	UploadedAt  int64  `json:"uploadedAt"`
	OwnerID     string `json:"ownerId"`
	Thumbnail   string `json:"thumbnail"`
}

func NewElasticVideo(client *esv7.Client) *Video {
	return &Video{
		client: client,
		index:  "video",
	}
}

func (v *Video) Index(ctx context.Context, video datastruct.Video) error {
	defer newOtelSpan(ctx, "Video.Index").End()

	b := indexedVideo{
		ID:          video.ID,
		Description: video.Description,
		Path:        video.Path,
		Thumbnail:   video.Thumbnail,
		Title:       video.Title,
		UploadedAt:  video.UploadedAt,
		OwnerID:     video.OwnerID,
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(b); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.IndexRequest{
		Index:      v.index,
		Body:       &buf,
		DocumentID: b.ID,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, v.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
	}

	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "IndexRequest.Do %d", resp.StatusCode)
	}

	io.Copy(io.Discard, req.Body)

	return nil

}

func (v *Video) Delete(ctx context.Context, id string) error {
	defer newOtelSpan(ctx, "Task.Delete").End()

	req := esv7api.DeleteRequest{
		Index:      v.index,
		DocumentID: id,
	}

	resp, err := req.Do(ctx, v.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "DeleteRequest.Do")
	}
	defer resp.Body.Close()

	if resp.IsError() {
		return internal.NewErrorf(internal.ErrorCodeUnknown, "DeleteRequest.Do %d", resp.StatusCode)
	}

	io.Copy(io.Discard, resp.Body)

	return nil
}
