package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Message struct {
	Name string `gorm:"not null"`
	Title string `gorm:"not null"`
	Sender User `gorm:"foreignKey:UserID"`
	UserID uint  `json:"sernder_id"`
	Description string `gorm:"not null"`
	Read bool
	gorm.Model
}
func (message Message) Validate() *httperors.HttpError{ 
	if message.Name == "" && len(message.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if message.Title == "" && len(message.Title) < 3 {
		return httperors.NewNotFoundError("Invalid Title")
	}
	
	if message.Description == "" && len(message.Description) < 10 {
		return httperors.NewNotFoundError("Invalid description")
	}
	return nil
}