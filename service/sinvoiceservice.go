package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	Sinvoiceservice sinvoiceservice = sinvoiceservice{}
)

type sinvoiceservice struct {
}

func (service sinvoiceservice) Create(sinvoice *model.SInvoice) (*httperors.HttpSuccess, *httperors.HttpError) {
	s, f := r.Sinvoicerepo.Create(sinvoice)
	if f != nil {
		return nil, f
	}
	return s, nil

}
func (service sinvoiceservice) GetOne(id int) (*model.SInvoice, *httperors.HttpError) {
	sinvoice, err1 := r.Sinvoicerepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return sinvoice, nil
}

func (service sinvoiceservice) GetAll(sinvoices []model.SInvoice, search *support.Search) ([]model.SInvoice, *httperors.HttpError) {
	sinvoices, err := r.Sinvoicerepo.GetAll(sinvoices, search)
	if err != nil {
		return nil, err
	}
	return sinvoices, nil
}
func (service sinvoiceservice) View() (*model.Soptions, *httperors.HttpError) {
	code, err1 := r.Sinvoicerepo.View()
	if err1 != nil {
		return nil, err1
	}
	return code, nil
}
// func (service sinvoiceservice) Update(id int, sinvoice *model.SInvoice) (*model.SInvoice, *httperors.HttpError) {
// 	sinvoice, err1 := r.Sinvoicerepo.Update(id, sinvoice)
// 	if err1 != nil {
// 		return nil, err1
// 	}

// 	return sinvoice, nil
// }
// func (service sinvoiceservice) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {

// 	success, failure := r.Sinvoicerepo.Delete(id)
// 	return success, failure
// }
