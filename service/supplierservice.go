package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/support"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
)

var (
	Supplierservice supplierservice = supplierservice{}

) 
type supplierservice struct {
	
}

func (service supplierservice) Create(supplier *model.Supplier) (string, *httperors.HttpError) {
	if err := supplier.Validate(); err != nil {
		return "", err
	}	
	s, err1 := r.Supplierrepo.Create(supplier)
	if err1 != nil {
		return "", err1
	}
	 return s, nil
 
}
func (service supplierservice) Login(asupplier *model.Loginsupplier) (*model.SupplierAuth, *httperors.HttpError) {
	
	supplier, err1 := r.Supplierrepo.Login(asupplier)
	if err1 != nil {
		return nil, err1
	} 
	return supplier, nil
}
func (service supplierservice) Forgot(email string) (string, *httperors.HttpError) {
	
	s, err1 := r.Supplierrepo.Forgot(email)
	if err1 != nil {
		return "", err1
	} 
	return s, nil
}
func (service supplierservice) Logout(token string) (*httperors.HttpError) {
	err1 := r.Supplierrepo.Logout(token)
	if err1 != nil {
		return err1
	}
	return nil
}
func (service supplierservice) GetOne(id int) (*model.Supplier, *httperors.HttpError) {
	supplier, err1 := r.Supplierrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return supplier, nil
}

func (service supplierservice) GetAll(search *support.Search) ([]model.Supplier, *httperors.HttpError) {
	
	results, err := r.Supplierrepo.GetAll(search)
	if err != nil { 
		return nil, err
	}
	return results, nil 
}

func (service supplierservice) Update(id int, supplier *model.Supplier) (*model.Supplier, *httperors.HttpError) {
	supplier, err1 := r.Supplierrepo.Update(id, supplier)
	if err1 != nil {
		return nil, err1
	}
	
	return supplier, nil
}
func (service supplierservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Supplierrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&supplier{})