package response

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type Login struct {
	Username    string `json:"username"`
	Name        string `json:"name"`
	AccessToken string `json:"access_token"`
}
