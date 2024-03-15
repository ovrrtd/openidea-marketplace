package entity

// Product represents a product entity in the database
type Product struct {
	ID            int64
	Name          string
	Price         int
	ImageURL      string
	Stock         int
	Condition     string
	Tags          string
	IsPurchasable bool
	PurchaseCount int
	UserID        int64
	User          User
	CreatedAt     int64
	UpdatedAt     int64
}

type GetAllProductFilter struct {
	UserOnly       bool
	UserID         int64
	Limit          int
	Offset         int
	Tags           []string
	Condition      string
	ShowEmptyStock bool
	MaxPrice       float64
	MinPrice       float64
	SortBy         string
	OrderBy        string
	Search         string
}
