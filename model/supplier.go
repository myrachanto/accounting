package model

import (
	
	"github.com/jinzhu/gorm"
	// "time"
	"golang.org/x/crypto/bcrypt"
	jwt "github.com/dgrijalva/jwt-go"
	"regexp"
	"github.com/myrachanto/accounting/httperors"
)

type Supplier struct {
	Name string `gorm:"not null" json:"name"`
	Company string `gorm:"not null" json:"company"`
	Phone string `gorm:"not null" json:"phone"`
	Address string `gorm:"not null" json:"address"`
	Picture string `json:"picture"`
	Email string `gorm:"not null;unique" json:"email"`
	Password string `gorm:"not null"`
	Invoices []Invoice `gorm:"foreignKey:supplierID; not null"`
	SupplierID uint `json:"supplierid"`
	gorm.Model
}
//Loginsupplier k
type Loginsupplier struct {
	Email string `gorm:"not null"`
	Password string `gorm:"not null"`
}
//CustomnerAuth str
type SupplierAuth struct {
	//supplier supplier `gorm:"foreignKey:supplierID; not null"`
	SupplierID uint `json:"supplierid"`
	Name string `json:"name"`
	Token string `gorm:"size:500;not null"`
	gorm.Model
}
//supplierToken struct declaration
type SupplierToken struct {
	SupplierID uint
	Name string `json:"name"`
	Email  string
	*jwt.StandardClaims
}
func (supplier Supplier)ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (supplier Supplier)ValidatePassword(password string) (bool, *httperors.HttpError) {
	if len(password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
func (supplier Supplier)HashPassword(password string)(string, *httperors.HttpError){
	pass, err := bcrypt.GenerateFromPassword([]byte(supplier.Password), 10)
		if err != nil {
			return "", httperors.NewNotFoundError("type a stronger password!")
		}
		return string(pass),nil 
		
	}
func (supplier Supplier) Compare(p1,p2 string) bool{
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	if err != nil {
		return false
	}
	return true
}
func (loginsupplier Loginsupplier) Validate() *httperors.HttpError{ 
	if loginsupplier.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if loginsupplier.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}
func (supplier Supplier) Validate() *httperors.HttpError{
	if supplier.Name == "" {
		return httperors.NewNotFoundError("Invalid first Name")
	}
	if supplier.Company == "" {
		return httperors.NewNotFoundError("Invalid last name")
	}
	if supplier.Phone == "" {
		return httperors.NewNotFoundError("Invalid phone number")
	}
	if supplier.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if supplier.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if supplier.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	// if supplier.Picture == "" {
	// 	return httperors.NewNotFoundError("Invalid picture")
	// }
	if supplier.Email == "" {
		return httperors.NewNotFoundError("Invalid picture")
	}
	return nil
}