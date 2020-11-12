package model

import (
	"github.com/jinzhu/gorm"
)
type Asstrans struct {
	Asset   Asset `gorm:"foreignKey:AssetID"`
	AssetID uint `json:"assetid"`
	Name string 
	Title string
	Description string 
	Payment Payment
	Depreciation float64
	Amount float64 
	Status bool 
	gorm.Model
}