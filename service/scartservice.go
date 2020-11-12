package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	// "github.com/myrachanto/accounting/support"
)

var (
	Scartservice scartservice = scartservice{}

) 
type scartservice struct {
	
}

func (service scartservice) Create(scart *model.Scart) (*model.Scart, *httperors.HttpError) {
	scart, err1 := r.Scartrepo.Create(scart)
	if err1 != nil {
		return nil, err1
	}
	 return scart, nil

}
func (service scartservice) GetOne(id int) (*model.Scart, *httperors.HttpError) {
	scart, err1 := r.Scartrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return scart, nil
}
func (service scartservice) View(id int) (*model.Options, *httperors.HttpError) {
	options, err1 := r.Scartrepo.View(id)
	if err1 != nil {
		return nil, err1
	}
	return options, nil
}
func (service scartservice) GetAll(scarts []model.Scart) ([]model.Scart, *httperors.HttpError) {
	scarts, err := r.Scartrepo.GetAll(scarts)
	if err != nil {
		return nil, err
	}
	return scarts, nil
}

func (service scartservice) Update(id int, scart *model.Scart) (*model.Scart, *httperors.HttpError) {
	scart, err1 := r.Scartrepo.Update(id, scart)
	if err1 != nil {
		return nil, err1
	}
	
	return scart, nil
}
func (service scartservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Scartrepo.Delete(id)
		return success, failure
}
///////deleting a batch////////////////////

//db.Where("age = ?", 20).Delete(&User{})