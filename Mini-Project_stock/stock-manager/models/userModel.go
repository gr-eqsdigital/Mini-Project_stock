package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fname    string `json:"fname"`
	Lname    string `json:"lname"`
	Phone    string `json:"phone"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}
