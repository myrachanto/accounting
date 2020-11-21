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
	Paymentformrepo paymentformrepo = paymentformrepo{}
)

///curtesy to gorm
type paymentformrepo struct{}

func (paymentformRepo paymentformrepo) Create(paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	if err := paymentform.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&paymentform)
	IndexRepo.DbClose(GormDB)
	return paymentform, nil
}
func (paymentformRepo paymentformrepo) GetOne(id int) (*model.Paymentform, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("paymentform with that id does not exists!")
	}
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&paymentform).Where("id = ?", id).First(&paymentform)
	IndexRepo.DbClose(GormDB)
	
	return &paymentform, nil
}

func (paymentformRepo paymentformrepo) GetAll(paymentforms []model.Paymentform,search *support.Search) ([]model.Paymentform, *httperors.HttpError) {

	results, err1 := paymentformRepo.Search(search, paymentforms)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (paymentformRepo paymentformrepo) All() (t []model.Paymentform, r *httperors.HttpError) {

	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&paymentform).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (paymentformRepo paymentformrepo) Update(id int, paymentform *model.Paymentform) (*model.Paymentform, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("paymentform with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	apaymentform := model.Paymentform{}
	
	GormDB.Model(&apaymentform).Where("id = ?", id).First(&apaymentform)
	if paymentform.Name  == "" {
		paymentform.Name = apaymentform.Name
	}
	if paymentform.Title  == "" {
		paymentform.Title = apaymentform.Title
	}
	if paymentform.Description  == "" {
		paymentform.Description = apaymentform.Description
	}
	fmt.Println(paymentform)
	GormDB.Model(&apaymentform).Where("id = ?", id).Update(&paymentform)
	
	IndexRepo.DbClose(GormDB)

	return paymentform, nil
}
func (paymentformRepo paymentformrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := paymentformRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&paymentform).Where("id = ?", id).First(&paymentform)
	GormDB.Delete(paymentform)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (paymentformRepo paymentformrepo)ProductUserExistByid(id int) bool {
	paymentform := model.Paymentform{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&paymentform, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (paymentformRepo paymentformrepo) Search(Ser *support.Search, paymentforms []model.Paymentform)([]model.Paymentform, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	paymentform := model.Paymentform{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&paymentform).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		fmt.Println(p.Results(&paymentforms))
				if err3 := p.Results(&paymentforms); err3 != nil {
					return nil, httperors.NewNotFoundError("something went wrong paginating!")
				}

		}else {

			q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
			p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
			p.SetPage(Ser.Page)
			fmt.Println(p.Results(&paymentforms,))
			if err3 := p.Results(&paymentforms); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&paymentforms);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&paymentforms); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return paymentforms, nil
}