package datastruct

type User struct {
	Username    string `json:"username"`
	Created_at  int64  `json:"created_at"`
	Followers   int64  `json:"followers"`
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
}
