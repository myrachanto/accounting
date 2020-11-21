package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Stransactionrepo ...
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
	GormDB.Save(&stransaction)
	
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
	res := GormDB.First(&stransaction, "id =?", id)
	if res.Error != nil {
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
		GormDB.Model(&stransaction).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&stransactions);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return stransactions, nil
}