package response

type Bank struct {
	ID            string `json:"bankAccountId"`
	Name          string `json:"bankName" validate:"required,min=5,max=15"`
	AccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
	AccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=15"`
	UserID        int64  `json:"userId"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}
