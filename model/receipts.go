package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Receipt struct {
	Name string `gorm:"not null"`
	Description string `gorm:"not null"` 
	Company string 
	CustomerID uint `gorm:"not null"`
	Customer Customer `gorm:"foreignKey:UserID; not null"`
	InvoiceID uint `gorm:"not null"`
	Invoice []Invoice `gorm:"foreignKey:UserID; not null"`
	PaymentMethod string `gorm:"not null"`
	Status string `gorm:"not null"`
	gorm.Model
}
func (receipts Receipt) Validate() *httperors.HttpError{ 
	if receipts.Name == "" && len(receipts.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if receipts.Description == "" && len(receipts.Description) > 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if receipts.CustomerID == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	if receipts.InvoiceID == 0{
		return httperors.NewNotFoundError("Invalid Invoice")
	}
	if receipts.PaymentMethod == "" {
		return httperors.NewNotFoundError("Invalid Payment method")
	}
	return nil
}