package response

type User struct {
	ID        int64  `json:"userId"`
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

type SellerDetail struct {
	ID               string `json:"userId"`
	Username         string `json:"username"`
	Name             string `json:"name"`
	ProductSoldTotal int    `json:"productSoldTotal"`
	Banks            []Bank `json:"bankAccounts"`
	CreatedAt        int64  `json:"created_at"`
	UpdatedAt        int64  `json:"updated_at"`
}
