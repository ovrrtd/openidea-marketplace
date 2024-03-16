package request

type CreateBank struct {
	Name          string `json:"bankName" validate:"required,min=5,max=15"`
	AccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
	AccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=15"`
	UserID        int64
}

type UpdateBank struct {
	ID            int64  `json:"id" validate:"required"`
	Name          string `json:"bankName" validate:"required,min=5,max=15"`
	AccountName   string `json:"bankAccountName" validate:"required,min=5,max=15"`
	AccountNumber string `json:"bankAccountNumber" validate:"required,min=5,max=15"`
	UserID        int64
}
