package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)


type Product struct {
	Name string `gorm:"not null" json:"name"`
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	//Subcategory Subcategory `gorm:"foreignKey:UserID; not null"`
	Category string ` json:"category"`
	Majorcategory string ` json:"majorcategory"`
	Picture string `json:"picture"`
	Price float64 `json:"price"`
	gorm.Model
}

type Info struct {
	Category []Category `gorm:"foreignKey:ategories; not null"`
	Majorcategory []Majorcategory `gorm:"foreignKey:Majorcategoryies; not null"`
	Subcategory []Subcategory `gorm:"foreignKey:Subcategories; not null"`
}
type Options struct {
	Price []Price `gorm:"foreignKey:prices; not null"`
	Tax []Tax `gorm:"foreignKey:taxs; not null"`
	Discount []Discount `gorm:"foreignKey:Discounts; not null"`
}
func (product Product) Validate() *httperors.HttpError{ 
	if product.Name == "" && len(product.Name) > 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if product.Title == "" && len(product.Title) > 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if product.Description == "" && len(product.Description) > 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}