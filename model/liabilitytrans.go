package model

import (
	"github.com/jinzhu/gorm"
)
type Liatran struct {
	Name string 
	Liability Liability `gorm:"foreignKey:LiabilityID"`
	LiabilityID uint  `json:"liabilityID"`
	Title string
	Description string
	Payment Payment
	Interest float64
	Amountpaid float64  
	Balance float64
	Status bool 
	gorm.Model
}