package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Price struct {
	Name string `gorm:"not null"`
	Product Product `gorm:"foreignKey:UserID; not null"`
	ProductID uint `json:"userid"`
	Amount float64 
	Buy bool
	gorm.Model
}
func (price Price) Validate() *httperors.HttpError{ 
	if price.Name == ""  {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if price.Amount == 0 {
		return httperors.NewNotFoundError("Invalid amount")
	}
	return nil
}