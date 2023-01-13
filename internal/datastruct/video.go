package datastruct

type Video struct {
	ID          string `json:"id"`
	Path        string `json:"path"`
	UploadedAt  int64  `json:"uploadedAt"`
	OwnerID     string `json:"ownerId"`
	Thumbnail   string `json:"thumbnail"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
