package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Tax struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	gorm.Model
}
func (tax Tax) Validate() *httperors.HttpError{ 
	if tax.Name == "" && len(tax.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if tax.Title == "" && len(tax.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if tax.Description == "" && len(tax.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}