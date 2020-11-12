package repository

import (
	"fmt"
	"strings"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
)

var (
	Paymentrepo paymentrepo = paymentrepo{}
)

///curtesy to gorm
type paymentrepo struct{}

func (paymentRepo paymentrepo) Create(payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&payment)
	IndexRepo.DbClose(GormDB)
	return payment, nil
}
func (paymentRepo paymentrepo) GetOne(id int) (*model.Payment, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Payrectrasans").Model(&payment).Where("id = ?", id).First(&payment)
	IndexRepo.DbClose(GormDB)
	
	return &payment, nil
}

func (paymentRepo paymentrepo) GetAll(payments []model.Payment,search *support.Search) ([]model.Payment, *httperors.HttpError) {
	
	results, err1 := paymentRepo.Search(search, payments)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (paymentRepo paymentrepo) Update(id int, payment *model.Payment) (*model.Payment, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	apayment := model.Payment{}
	
	GormDB.Model(&apayment).Where("id = ?", id).First(&apayment)
	// if payment.payment  == "" {
	// 	payment.payment = apayment.payment
	// }
	// if payment.Description  == "" {
	// 	payment.Description = apayment.Description
	// }
	// if payment.Subtotal  == 0 {
	// 	payment.Subtotal = apayment.Subtotal
	// }
	// if payment.Discount  == 0 {
	// 	payment.Discount = apayment.Discount
	// }	
	// if payment.AmountPaid  == 0 {
	// 	payment.AmountPaid = apayment.AmountPaid
	// }
	GormDB.Model(&payment).Where("id = ?", id).First(&payment).Update(&apayment)
	
	IndexRepo.DbClose(GormDB)

	return payment, nil
}
func (paymentRepo paymentrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := paymentRepo.paymentUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("payment with that id does not exists!")
	}
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&payment).Where("id = ?", id).First(&payment)
	GormDB.Delete(payment)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (paymentRepo paymentrepo)paymentUserExistByid(id int) bool {
	payment := model.Payment{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&payment, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (paymentRepo paymentrepo) Search(Ser *support.Search, payments []model.Payment)([]model.Payment, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	payment := model.Payment{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Preload("Payrectrasans").Model(&payment).Order(Ser.Column+" "+Ser.Direction).Find(&payments)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&payments);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Preload("Payrectrasans").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Preload("Payrectrasans").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&payments);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&payments); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return payments, nil
}