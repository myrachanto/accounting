package model

import (
	
	"github.com/jinzhu/gorm"
	// "time"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/myrachanto/accounting/httperors"
)

type Customer struct {
	Name string `gorm:"not null" json:"name"`
	Company string `gorm:"not null" json:"company"`
	Phone string `gorm:"not null" json:"phone"`
	Address string `gorm:"not null" json:"address"`
	Picture string `json:"picture"`
	Email string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null"`
	Invoices []Invoice `gorm:"foreignKey:CustomerID; not null"`
	CustomerID uint `json:"customerid"`
	gorm.Model
}
//Logincustomer k
type Logincustomer struct {
	Email string `gorm:"not null"`
	Password string `gorm:"not null"`
}
//CustomnerAuth str
type CustomnerAuth struct {
	//customer customer `gorm:"foreignKey:customerID; not null"`
	CustomerID uint `json:"customerid"`
	Name string `json:"name"`
	Token string `gorm:"size:500;not null"`
	gorm.Model
}
//CustomerToken struct declaration
type CustomerToken struct {
	CustomerID uint
	Name string `json:"name"`
	Email  string
	*jwt.StandardClaims
}
func (customer Customer)ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (customer Customer)ValidatePassword(password string) (bool, *httperors.HttpError) {
	if len(password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
func (customer Customer)HashPassword(password string)(string, *httperors.HttpError){
	pass, err := bcrypt.GenerateFromPassword([]byte(customer.Password), 10)
		if err != nil {
			return "", httperors.NewNotFoundError("type a stronger password!")
		}
		return string(pass),nil 
		
	}
func (customer Customer) Compare(p1,p2 string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	if err != nil {
		return false
	}
	return true
}
func (logincustomer Logincustomer) Validate() *httperors.HttpError{ 
	if logincustomer.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if logincustomer.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}
func (customer Customer) Validate() *httperors.HttpError{
	if customer.Name == "" {
		return httperors.NewNotFoundError("Invalid first Name")
	}
	if customer.Company == "" {
		return httperors.NewNotFoundError("Invalid last name")
	}
	if customer.Phone == "" {
		return httperors.NewNotFoundError("Invalid phone number")
	}
	if customer.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if customer.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if customer.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	// if customer.Picture == "" {
	// 	return httperors.NewNotFoundError("Invalid picture")
	// }
	if customer.Email == "" {
		return httperors.NewNotFoundError("Invalid picture")
	}
	return nil
}