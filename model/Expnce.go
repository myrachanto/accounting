package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Expence struct {
	Name string `gorm:"not null"`
	Description string `gorm:"not null"`
	Company string 
	Amount float64
	Paid string `gorm:"not null"`
	gorm.Model
}
func (expence Expence) Validate() *httperors.HttpError{ 
	if expence.Name == "" && len(expence.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if expence.Description == "" && len(expence.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if expence.Company == "" {
		return httperors.NewNotFoundError("Invalid description")
	}
	if expence.Amount == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	return nil
}