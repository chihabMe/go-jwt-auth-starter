// @Title
// @Description
// @Author
// @Update
package models

type User struct {
	Base
	Username string `json:"username" validate:"required,min=5,max=20"`
	Email    string `json:"email" gorm:"unique" validate:"required,email"`
	Password string `json:"password" gorm:"unique" validate:"reqiured,min=8"`
	Token    string `json:"token"`
	Refresh  string `json:"refresh"`
	UserType string `json:"type" validate:"eq=ADMIN|eq=USER"`
}
