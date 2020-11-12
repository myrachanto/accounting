package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Expencetrasan struct {
	Name string 
	Expence Expence `gorm:"foreignKey:ExpenceID"`
	ExpenceID uint  `json:"expenceID"`
	Title string
	Description string
	Amount float64 
	Status bool 
	Paid string `gorm:"not null"`
	gorm.Model
}
func (expence Expencetrasan) Validate() *httperors.HttpError{ 
	if expence.Name == "" && len(expence.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if expence.Description == "" && len(expence.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if expence.Amount == 0 {
		return httperors.NewNotFoundError("Invalid customer")
	}
	return nil
}