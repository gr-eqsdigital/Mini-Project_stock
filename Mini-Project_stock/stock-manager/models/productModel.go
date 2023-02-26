package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string
	Description string
	Quantity    int
	Price       float64
	CategoryID  uint
}
