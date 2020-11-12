package model

import (
	"github.com/jinzhu/gorm"
)
type Payrectrasan struct {
	Name string  
	Title string
	Customerid uint
	CustomerName uint
	Receipt Receipt `gorm:"foreignKey:ReceiptID"`
	ReceiptID uint  `json:"receiptID"`
	Supplierid uint
	SupplierName uint
	Payment Payment `gorm:"foreignKey:PaymentID"`
	PaymentID uint  `json:"paymentID"`
	Description string
	Amount float64 
	Status bool 
	gorm.Model
}