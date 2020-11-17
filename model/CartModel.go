package model

import (
	"github.com/jinzhu/gorm"
)
type Cart struct {
	Product   Product `gorm:"foreignKey:productid"`
	ProductID uint `json:"productid"`
	Code string `json:"code" json:"code"`
	Name string `gorm:"not null" json:"name"`
	Quantity float64 `gorm:"not null" json:"quantity"`
	Price float64 `gorm:"not null" son:"price"`
	Subtotal float64 `gorm:"not null"`
	Discount float64 `json:"discount"`
	Tax float64 `json:"tax"`
	Total float64 
	Cartstatus bool 
	Picture string 
	gorm.Model
}
type Scart struct {
	Product   Product `gorm:"foreignKey:productid"`
	ProductID uint `json:"productid"`
	Code string
	Name string `gorm:"not null"`
	Qty float64 `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Subtotal float64 `gorm:"not null"`
	Discount float64 
	Tax float64 
	Total float64 
	Cartstatus bool  
	Picture string 
	gorm.Model
}