package dto

type VideoSearchParams struct {
	Title *string `json:"title"`
	From  int64   `json:"from"`
	Size  int64   `json:"size"`
}
