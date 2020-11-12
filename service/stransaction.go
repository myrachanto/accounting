package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	Stransactionservice stransactionservice = stransactionservice{}

) 
type stransactionservice struct {
	
}

func (service stransactionservice) Create(stransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	stransaction, err1 := r.Stransactionrepo.Create(stransaction)
	if err1 != nil {
		return nil, err1
	}
	 return stransaction, nil

}
func (service stransactionservice) GetOne(id int) (*model.STransaction, *httperors.HttpError) {
	stransaction, err1 := r.Stransactionrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return stransaction, nil
}

func (service stransactionservice) GetAll(stransactions []model.STransaction,search *support.Search) ([]model.STransaction, *httperors.HttpError) {
	stransactions, err := r.Stransactionrepo.GetAll(stransactions,search)
	if err != nil {
		return nil, err
	}
	return stransactions, nil
}

func (service stransactionservice) Update(id int, stransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	stransaction, err1 := r.Stransactionrepo.Update(id, stransaction)
	if err1 != nil {
		return nil, err1
	}
	
	return stransaction, nil
}
func (service stransactionservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Stransactionrepo.Delete(id)
		return success, failure
}
