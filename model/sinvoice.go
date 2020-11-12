package model

import (
	"time"
	"github.com/jinzhu/gorm"
)
type SInvoice struct {
	Supplier Supplier `gorm:"foreignKey:SupplierID; not null"`
	SupplierID uint `json:"supplier"`
	Code string
	Title  string
	Description string
	Scart []Scart `gorm:"foreignKey:ScartID; not null"`
	ScartID uint `json:"scart"`
	Dated time.Time
	Due_date *time.Time
	Sub_total float64 
	Discount float64 
	Tax float64 
	Total float64 
	PaidStatus bool 
	AmountPaid float64 
	Balance float64
	Status bool
	Cn bool
	STransaction []STransaction `gorm:"foreignKey:STransactionID; not null"`
	STransactionID uint `json:"stransaction"`
	gorm.Model
}
type Soptions struct {
	Code string
	Supplier []Supplier
	Product []Product
}
type Poptions struct {
	Supplier []Supplier
}