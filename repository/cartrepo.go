package repository

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
)
//Cartrepo...
var (
	Cartrepo cartrepo = cartrepo{}
)
//Totals ...
type Totals struct {
	Discount float64
	Subtotal float64
	Total float64
}
///curtesy to gorm
type cartrepo struct{}
//////////////
////////////TODO user id///////////
/////////////////////////////////////////
func (cartRepo cartrepo) Create(cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// code, err := Invoicerepo.GeneCode()
	// if err != nil {
	// 	return nil, err
	// }
	// cart.Code = code
	
	grossamount := cart.Quantity * cart.SPrice
	taxamount := cart.Tax/100 * grossamount
	discountamount := cart.Discount/100 * grossamount
	// fmt.Println(grossamount,taxamount,discountamount)
	// fmt.Println(cart)
	cart.Total = grossamount - discountamount + taxamount
	code := cart.Code
	cart.Tax = taxamount
	cart.Discount = discountamount
	ok := Invoicerepo.InvoiceExistByCode(code)
	if ok == true {
		return nil, httperors.NewNotFoundError("That invoice is already saved!")
	}
	GormDB.Create(&cart)
	IndexRepo.DbClose(GormDB)
	return cart, nil
}
func (cartRepo cartrepo) View(code string) ([]model.Cart, *httperors.HttpError) {
	mc, e := cartRepo.Getcarts(code)
	if e != nil{
		return nil, e
	}
	return mc, nil
}
func (cartRepo cartrepo) Getcarts(code string) (t []model.Cart, e *httperors.HttpError) {

	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&cart).Where("code = ?", code).Find(&t)
	IndexRepo.DbClose(GormDB)
	
	return t, nil
}
func (cartRepo cartrepo) GetOne(id int) (*model.Cart, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&cart).Where("id = ?", id).First(&cart)
	IndexRepo.DbClose(GormDB)
	
	return &cart, nil
}
func (cartRepo cartrepo) All() (t []model.Cart, r *httperors.HttpError) {

	cart := model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&cart).Order("name").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (cartRepo cartrepo) GetAll(carts []model.Cart) ([]model.Cart, *httperors.HttpError) {
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cart := model.Cart{}
	GormDB.Model(&cart).Find(&carts)
	
	IndexRepo.DbClose(GormDB)
	if len(carts) == 0 {
		return nil, httperors.NewNotFoundError("No results found!")
	}
	return carts, nil
}

func (cartRepo cartrepo) Update(id int, cart *model.Cart) (*model.Cart, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	acart := model.Cart{}
	
	GormDB.Model(&acart).Where("id = ?", id).First(&acart)
	if cart.Name  == "" {
		cart.Name = acart.Name
	}
	if cart.Quantity  == 0 {
		cart.Quantity = acart.Quantity
	}
	if cart.SPrice  == 0 {
		cart.SPrice = acart.SPrice
	}
	
	if cart.Discount  == 0 {
		cart.Discount = acart.Discount
	}
	if cart.Tax  == 0 {
		cart.Tax = acart.Tax
	}
	GormDB.Save(&cart)
	
	IndexRepo.DbClose(GormDB)

	return cart, nil
}
func (cartRepo cartrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := cartRepo.cartUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("cart with that id does not exists!")
	}
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&cart).Where("id = ?", id).First(&cart)
	GormDB.Delete(cart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo) DeleteAll(code string) (*httperors.HttpSuccess, *httperors.HttpError) {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Where("code = ?", code).Delete(&cart)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (cartRepo cartrepo)cartUserExistByid(id int) bool {
	cart := model.Cart{}
	GormDB, err1 :=IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&cart, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (cartRepo cartrepo)SumTotal(code string) (Total *Totals, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)
	for _, t := range carts {
		Total.Discount += t.Discount
		Total.Subtotal += t.Subtotal
		Total.Total += t.Total
	}
	return Total,nil
}
func (cartRepo cartrepo)CarttoTransaction(code string) (tr []model.Transaction, err *httperors.HttpError) {
	carts := []model.Cart{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil,httperors.NewNotFoundError("db connection failed!")
	}
	GormDB.Where("code = ?", code).Find(&carts)
	IndexRepo.DbClose(GormDB)
	for _, c := range carts {
		trans := model.Transaction{Productname:c.Name,Productid:c.ProductID,Quantity: c.Quantity,Price: c.SPrice,Tax:c.Tax, Code:code, Subtotal:c.Subtotal, Discount:c.Discount,Total:c.Total}
		tr = append(tr, trans)
	}
	return tr,nil
}
