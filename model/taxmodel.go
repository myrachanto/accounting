package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Tax struct {
	Name string `gorm:"not null"`
	Product Product `gorm:"foreignKey:UserID; not null"`
	ProductID uint `json:"userid"`
	Amount float64 
	gorm.Model
}
func (tax Tax) Validate() *httperors.HttpError{ 
	if tax.Name == ""  {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if tax.Amount == 0 {
		return httperors.NewNotFoundError("Invalid amount")
	}
	return nil
}