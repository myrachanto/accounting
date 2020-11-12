package repository

import (
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	// "github.com/myrachanto/accounting/support"
)

var (
	Scartrepo scartrepo = scartrepo{}
)

///curtesy to gorm
type scartrepo struct{}

func (scartRepo scartrepo) Create(scart *model.Scart) (*model.Scart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	code := scart.Code
	scart.Subtotal = scart.Price * scart.Qty 
	scart.Total = scart.Subtotal - scart.Discount
	ok := Invoicerepo.InvoiceExistByCode(code)
	if ok == true {
		return nil, httperors.NewNotFoundError("That invoice is already saved!")
	}
	GormDB.Create(&scart)
	IndexRepo.DbClose(GormDB)
	return scart, nil
}
func (scartRepo scartrepo) GetOne(id int) (*model.Scart, *httperors.HttpError) {
	ok := scartRepo.scartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("scart with that id does not exists!")
	}
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&scart).Where("id = ?", id).First(&scart)
	IndexRepo.DbClose(GormDB)
	
	return &scart, nil
}
func (scartRepo scartrepo) View(id int) (*model.Options, *httperors.HttpError) {
	c,e := Pricerepo.GetOptionsell(id);if e != nil {
		return nil,e
	}
	m,me := Taxrepo.GetOption();if me != nil {
		return nil,me
	}
	s,se := Discountrepo.GetOptionsell(id);if se != nil {
		return nil,se
	}
	options := model.Options{}
	options.Price = c
	options.Tax = m
	options.Discount = s
	return &options, nil
}
func (scartRepo scartrepo) GetAll(scarts []model.Scart) ([]model.Scart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	scart := model.Scart{}
	GormDB.Model(&scart).Find(&scarts)
	
	IndexRepo.DbClose(GormDB)
	if len(scarts) == 0 {
		return nil, httperors.NewNotFoundError("No results found!")
	}
	return scarts, nil
}

func (scartRepo scartrepo) Update(id int, scart *model.Scart) (*model.Scart, *httperors.HttpError) {
	ok := scartRepo.scartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("scart with that id does not exists!")
	}
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	ascart := model.Scart{}
	
	GormDB.Model(&ascart).Where("id = ?", id).First(&ascart)
	if scart.Name  == "" {
		scart.Name = ascart.Name
	}
	if scart.Qty  == 0 {
		scart.Qty = ascart.Qty
	}
	if scart.Price  == 0 {
		scart.Price = ascart.Price
	}
	
	if scart.Discount  == 0 {
		scart.Discount = ascart.Discount
	}
	if scart.Tax  == 0 {
		scart.Tax = ascart.Tax
	}
	GormDB.Model(&scart).Where("id = ?", id).First(&scart).Update(&ascart)
	
	IndexRepo.DbClose(GormDB)

	return scart, nil
}
func (scartRepo scartrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := scartRepo.scartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("scart with that id does not exists!")
	}
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&scart).Where("id = ?", id).First(&scart)
	GormDB.Delete(scart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (scartRepo scartrepo)scartUserExistByid(id int) bool {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&scart, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (scartRepo scartrepo) DeleteAll(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	scart := model.Scart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Where("code = ?", code).Delete(&scart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}

func (scartRepo scartrepo)SumTotal(code string) (Total *Totals, err *httperors.HttpError) {
	scarts := []model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&scarts)
	IndexRepo.DbClose(GormDB)
	for _, t := range scarts {
		Total.Discount += t.Discount
		Total.Subtotal += t.Subtotal
		Total.Total += t.Total
	}
	return Total,nil
}
func (scartRepo scartrepo)CarttoTransaction(code string) (tr []model.STransaction, err *httperors.HttpError) {
	scarts := []model.Scart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&scarts)
	IndexRepo.DbClose(GormDB)
	for _, c := range scarts {
		trans := model.STransaction{Productname:c.Name,Productid:c.ProductID,Qty: c.Qty,Price: c.Price,Tax:c.Tax, Code:code, Subtotal:c.Subtotal, Discount:c.Discount,Total:c.Total}
		tr = append(tr, trans) 
	}
	return tr,nil
}

