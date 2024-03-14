package entity

// Product represents a product entity in the database
type Product struct {
	ID              int64
	Label           string
	LabelMask       string
	Description     string
	DescriptionMask string
	Price           int64
	UserID          int64
	CreatedAt       int64
	UpdatedAt       int64
}
