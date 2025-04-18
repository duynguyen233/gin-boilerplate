package model

type User struct {
	ID        int    `json:"id" gorm:"primaryKey,autoIncrement"`
	Username  string `json:"username" gorm:"unique"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Address   string `json:"address"`
	Phone     string `json:"phone"`
	Role      string `json:"role"`
}
