package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Payment struct {
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
func (payment Payment) Validate() *httperors.HttpError{ 
	if payment.Name == "" && len(payment.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if payment.Description == "" && len(payment.Description) > 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if payment.CustomerID == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	if payment.InvoiceID == 0{
		return httperors.NewNotFoundError("Invalid Invoice")
	}
	if payment.PaymentMethod == "" {
		return httperors.NewNotFoundError("Invalid Payment method")
	}
	return nil
}