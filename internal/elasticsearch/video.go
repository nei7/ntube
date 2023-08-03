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
	"github.com/nei7/ntube/internal/dto"
)

type Video struct {
	client *esv7.Client
	index  string
}

type indexedVideo struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Title       string          `json:"title"`
	Path        string          `json:"path"`
	UploadedAt  int64           `json:"uploadedAt"`
	User        datastruct.User `json:"user"`
	Thumbnail   string          `json:"thumbnail"`
}

func NewVideo(client *esv7.Client) *Video {
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
		User:        video.User,
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

func (v *Video) Search(ctx context.Context, params dto.VideoSearchParams) (datastruct.SearchResult, error) {

	defer newOtelSpan(ctx, "Task.Search").End()

	should := make([]interface{}, 0, 1)

	if params.Title != nil {
		should = append(should, map[string]interface{}{
			"match": map[string]interface{}{
				"title": *&params.Title,
			}})
	}

	var query map[string]interface{}

	if len(should) > 1 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"should": should,
				},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": should[0],
		}
	}

	// query["sort"] = []interface{}{
	// 	"_score",
	// 	map[string]interface{}{"id": "asc"},
	// }

	query["from"] = params.From
	query["size"] = params.Size

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return datastruct.SearchResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewEncoder.Encode")
	}

	req := esv7api.SearchRequest{
		Index: []string{v.index},
		Body:  &buf,
	}

	resp, err := req.Do(ctx, v.client)
	if err != nil {
		return datastruct.SearchResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do")
	}

	if resp.IsError() {
		return datastruct.SearchResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "SearchRequest.Do %d", resp.StatusCode)
	}

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source indexedVideo `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		return datastruct.SearchResult{}, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.NewDecoder.Decode")
	}

	res := make([]datastruct.Video, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		res[i] = datastruct.Video{
			ID:          hit.Source.ID,
			Description: hit.Source.Description,
			Title:       hit.Source.Title,
			Path:        hit.Source.Path,
			Thumbnail:   hit.Source.Thumbnail,
			User:        hit.Source.User,
		}
	}

	return datastruct.SearchResult{
		Videos: res,
		Total:  hits.Hits.Total.Value,
	}, nil
}
