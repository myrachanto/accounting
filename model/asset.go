package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Asset struct {
	Name string `gorm:"not null"`
	Description string `gorm:"not null"`
	Ownership string 
	Price float64 `gorm:"not null"`
	Depreciationtype string
	Depreciationrate float64
	ExpectedUsage float64
	Liscence string
	gorm.Model
}
func (asset Asset) Validate() *httperors.HttpError{ 
	if asset.Name == "" && len(asset.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if asset.Description == "" && len(asset.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if asset.Liscence == "" {
		return httperors.NewNotFoundError("Invalid liscence")
	}
	if asset.Depreciationtype == "" {
		return httperors.NewNotFoundError("Invalid depreciation type")
	}
	if asset.Depreciationrate < 0 {
		return httperors.NewNotFoundError("Invalid depreciation rate")
	}
	if asset.Price < 0 {
		return httperors.NewNotFoundError("Invalid Price")
	}
	return nil
}