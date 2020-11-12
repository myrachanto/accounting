package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)
type Nortification struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Sender User `gorm:"foreignKey:UserID"`
	UserID uint  `json:"sernder_id"`
	Description string `gorm:"not null"`
	Read bool
	gorm.Model
}
func (nortification Nortification) Validate() *httperors.HttpError{ 
	if nortification.Name == "" && len(nortification.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if nortification.Title == "" && len(nortification.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if nortification.Description == "" && len(nortification.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}