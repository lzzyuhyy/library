package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);not null;"`
	Phone    string `gorm:"type:char(11);not null;"`
}
