package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Discount struct {
	Name string `gorm:"not null"`
	Product Product `gorm:"foreignKey:UserID; not null"`
	ProductID uint `json:"userid"`
	Buy bool
	Amount float64 
	gorm.Model
}
func (discount Discount) Validate() *httperors.HttpError{ 
	if discount.Name == ""  {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if discount.Amount == 0 {
		return httperors.NewNotFoundError("Invalid amount")
	}
	return nil
}