package request

type Product struct {
	Label       string `json:"label" validate:"required,max=100"`
	Description string `json:"description" validate:"required,max=1000"`
	Price       int64  `json:"price" validate:"required,gte=0"`
	UserID      int64  `validate:"required"`
}
