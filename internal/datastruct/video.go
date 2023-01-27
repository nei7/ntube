package datastruct

type Video struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	UploadedAt  int64  `json:"uploadedAt"`
	User        User   `json:"user"`
	Thumbnail   string `json:"thumbnail"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type SearchResult struct {
	Videos []Video `json:"videos"`
	Total  int64   `json:"total"`
}
