package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname    string
	Lname    string
	Phone    string
	Email    string `gorm:"unique"`
	Password string
}
