package request

type Product struct {
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         int      `json:"price" validate:"required,min=0"`
	ImageURL      string   `json:"imageUrl" validate:"required,url"`
	Stock         int      `json:"stock" validate:"required,min=0"`
	Condition     string   `json:"condition" validate:"required"`
	Tags          []string `json:"tags" validate:"required,min=1,max=5"`
	IsPurchasable bool     `json:"isPurchasable"`
	UserID        int64
}

type UpdateProduct struct {
	ID            int64    `json:"id" validate:"required"`
	Name          string   `json:"name" validate:"required,min=5,max=60"`
	Price         int      `json:"price" validate:"required,min=0"`
	ImageURL      string   `json:"imageUrl" validate:"required,url"`
	Condition     string   `json:"condition" validate:"required"`
	Tags          []string `json:"tags" validate:"required,min=1,max=5"`
	IsPurchasable bool     `json:"isPurchasable"`
}

type UpdateProductStock struct {
	ID    int64
	Stock int `json:"stock" validate:"required,min=0"`
}

type GetProducts struct {
	UserID         int64
	UserOnly       bool     `query:"userOnly"`
	Limit          int      `query:"limit" default:"10"`
	Offset         int      `query:"offset" default:"0"`
	Tags           []string `query:"tags"`
	Condition      string   `query:"condition"`
	ShowEmptyStock bool     `query:"showEmptyStock"`
	MaxPrice       float64  `query:"maxPrice"`
	MinPrice       float64  `query:"minPrice"`
	SortBy         string   `query:"sortBy"`
	OrderBy        string   `query:"orderBy"`
	Search         string   `query:"search"`
}
