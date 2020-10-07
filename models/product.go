package models

import "gorm.io/gorm"

// Product is Gorm model of product
type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	IsAvailable bool    `json:"isAvailable"`
}
