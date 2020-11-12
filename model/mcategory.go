package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Majorcategory struct {
	Name string `gorm:"not null" json:"name"` 
	Title string `gorm:"not null" json:"title"`
	Description string `gorm:"not null" json:"description"`
	Category []Category
	gorm.Model
}
type MOptions struct{
	MajorCats []Majorcategory
}
func (majorcategory Majorcategory) Validate() *httperors.HttpError{ 
	if majorcategory.Name == "" && len(majorcategory.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if majorcategory.Title == "" && len(majorcategory.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if majorcategory.Description == "" && len(majorcategory.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}