package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
)

var (
	Stransactionrepo stransactionrepo = stransactionrepo{}
)

///curtesy to gorm
type stransactionrepo struct{}

func (stransactionRepo stransactionrepo) Create(stransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&stransaction)
	IndexRepo.DbClose(GormDB)
	return stransaction, nil
}
func (stransactionRepo stransactionrepo) GetOne(id int) (*model.STransaction, *httperors.HttpError) {
	ok := stransactionRepo.stransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("stransaction with that id does not exists!")
	}
	stransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&stransaction).Where("id = ?", id).First(&stransaction)
	IndexRepo.DbClose(GormDB)
	
	return &stransaction, nil
}

func (stransactionRepo stransactionrepo) GetAll(stransactions []model.STransaction,search *support.Search) ([]model.STransaction, *httperors.HttpError) {
	results, err1 := stransactionRepo.Search(search, stransactions)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (stransactionRepo stransactionrepo) Update(id int, stransaction *model.STransaction) (*model.STransaction, *httperors.HttpError) {
	ok := stransactionRepo.stransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("stransaction with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	astransaction := model.STransaction{}
	
	GormDB.Model(&astransaction).Where("id = ?", id).First(&astransaction)
	// if stransaction.Name  == "" {
	// 	stransaction.Name = astransaction.Name
	// }
	// if stransaction.Qty  == 0 {
	// 	stransaction.Qty = astransaction.Qty
	// }
	// if stransaction.Price  == 0 {
	// 	stransaction.Price = astransaction.Price
	// }
	
	// if stransaction.Discount  == 0 {
	// 	stransaction.Discount = astransaction.Discount
	// }
	// if stransaction.Tax  == 0 {
	// 	stransaction.Tax = astransaction.Tax
	// }
	GormDB.Model(&stransaction).Where("id = ?", id).First(&stransaction).Update(&astransaction)
	
	IndexRepo.DbClose(GormDB)

	return stransaction, nil
}
func (stransactionRepo stransactionrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := stransactionRepo.stransactionUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("stransaction with that id does not exists!")
	}
	stransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&stransaction).Where("id = ?", id).First(&stransaction)
	GormDB.Delete(stransaction)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (stransactionRepo stransactionrepo)stransactionUserExistByid(id int) bool {
	stransaction := model.STransaction{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&stransaction, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (stransactionRepo stransactionrepo) Search(Ser *support.Search, stransactions []model.STransaction)([]model.STransaction, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	stransaction := model.STransaction{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&stransaction).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&stransactions); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return stransactions, nil
}