package entity

type User struct {
	ID        int64
	Name      string
	Username  string
	Password  string
	Banks     []Bank
	CreatedAt int64
	UpdatedAt int64
}
