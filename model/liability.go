package model

import (
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/accounting/httperors"
)

type Liability struct {
	Name string `gorm:"not null"`
	Description string `gorm:"not null"`
	Creditor string 
	Approvedby string 
	Amount float64 `gorm:"not null"`
	Interestrate float64
	Paymentperiod float64
	Amoutinterest float64
	Monthlypayment float64
	gorm.Model
}
func (liability Liability) Validate() *httperors.HttpError{ 
	if liability.Name == "" && len(liability.Name) < 3 {
		return httperors.NewNotFoundError("Invalid Name")
	}
	if liability.Description == "" && len(liability.Description) < 3 {
		return httperors.NewNotFoundError("Invalid description")
	}
	if liability.Creditor == "" {
		return httperors.NewNotFoundError("Invalid Creditor name")
	}
	if liability.Approvedby == "" {
		return httperors.NewNotFoundError("Invalid Approved name")
	}
	if liability.Paymentperiod < 0 {
		return httperors.NewNotFoundError("Invalid payment period")
	}
	if liability.Interestrate < 0 {
		return httperors.NewNotFoundError("Invalid interst rate")
	}
	if liability.Amount < 0 {
		return httperors.NewNotFoundError("Invalid amount")
	}
	return nil
}