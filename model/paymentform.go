package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Paymentform struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	gorm.Model
}
func (paymentform Paymentform) Validate() *httperors.HttpError{ 
	if paymentform.Name == "" && len(paymentform.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if paymentform.Title == "" && len(paymentform.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if paymentform.Description == "" && len(paymentform.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}