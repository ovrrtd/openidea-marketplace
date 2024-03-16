package entity

type Bank struct {
	ID            int64
	Name          string
	AccountName   string
	AccountNumber string
	UserID        int64
	User          User
	CreatedAt     int64
	UpdatedAt     int64
}
