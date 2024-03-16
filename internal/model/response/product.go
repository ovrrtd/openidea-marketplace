package response

type Product struct {
	ID            string   `json:"productId"`
	Name          string   `json:"name"`
	Price         int      `json:"price"`
	ImageURL      string   `json:"imageUrl"`
	Stock         int      `json:"stock"`
	Condition     string   `json:"condition"`
	Tags          []string `json:"tags"`
	IsPurchasable bool     `json:"isPurchasable"`
	PurchaseCount int      `json:"purchaseCount"`
	UserID        int64    `json:"user_id"`
	CreatedAt     int64    `json:"created_at"`
	UpdatedAt     int64    `json:"updated_at"`
}

type PurchaseProduct struct {
	BankAccountId        string `json:"bankAccountId" validate:"required"`
	PaymentProofImageUrl string `json:"paymentProofImageUrl" validate:"required,url"`
	Quantity             int    `json:"quantity" validate:"required,min=1"`
}
