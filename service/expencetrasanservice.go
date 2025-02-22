package service

import (
	// "fmt"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	r "github.com/myrachanto/accounting/repository"
	"github.com/myrachanto/accounting/support"
)

var (
	ExpencetrasanService expencetrasanService = expencetrasanService{}

) 
type expencetrasanService struct {
	
}

func (service expencetrasanService) Create(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	if err := expencetrasan.Validate(); err != nil {
		return nil, err
	}	
	expencetrasan, err1 := r.Expencetrasanrepo.Create(expencetrasan)
	if err1 != nil {
		return nil, err1
	}
	 return expencetrasan, nil

}
func (service expencetrasanService) GetOne(id int) (*model.Expencetrasan, *httperors.HttpError) {
	expencetrasan, err1 := r.Expencetrasanrepo.GetOne(id)
	if err1 != nil {
		return nil, err1
	}
	return expencetrasan, nil
}

func (service expencetrasanService) GetAll(expencetrasans []model.Expencetrasan,search *support.Search) ([]model.Expencetrasan, *httperors.HttpError) {
	expencetrasans, err := r.Expencetrasanrepo.GetAll(expencetrasans,search)
	if err != nil {
		return nil, err
	}
	return expencetrasans, nil
}

func (service expencetrasanService) Update(id int, expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	expencetrasan, err1 := r.Expencetrasanrepo.Update(id, expencetrasan)
	if err1 != nil {
		return nil, err1
	}
	
	return expencetrasan, nil
}
func (service expencetrasanService) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	
		success, failure := r.Expencetrasanrepo.Delete(id)
		return success, failure
}
