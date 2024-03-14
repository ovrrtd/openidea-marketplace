package response

type Product struct {
	ID              int64  `json:"id"`
	Label           string `json:"label"`
	LabelMask       string `json:"label_mask"`
	Description     string `json:"description"`
	DescriptionMask string `json:"description_mask"`
	Price           int64  `json:"price"`
	UserID          int64  `json:"user_id"`
	CreatedAt       int64  `json:"created_at"`
	UpdatedAt       int64  `json:"updated_at"`
}
