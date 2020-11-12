package model

import (
	"github.com/jinzhu/gorm"
)
type STransaction struct {
	Code string
	Productname string 
	Productid uint
	Title string
	Qty float64 
	Price float64 
	Tax float64 
	Subtotal float64 
	Discount float64 
	Total float64 
	AmountPaid float64 
	Balance float64
	Status bool 
	SInvoice SInvoice `gorm:"foreignKey:SInvoiceID; not null"`
	SInvoiceID uint `json:"sinvoice"`
	gorm.Model
}