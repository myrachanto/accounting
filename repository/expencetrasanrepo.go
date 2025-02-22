package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Expencetrasanrepo ...
var (
	Expencetrasanrepo expencetrasanrepo = expencetrasanrepo{}
)

///curtesy to gorm
type expencetrasanrepo struct{}

func (expencetrasanRepo expencetrasanrepo) Create(expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	if err := expencetrasan.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	id := expencetrasan.Expence.ID
	expencetrasan.ExpenceID = id
	GormDB.Create(&expencetrasan)
	IndexRepo.DbClose(GormDB)
	return expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) GetOne(id int) (*model.Expencetrasan, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expencetrasan with that id does not exists!")
	}
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&expencetrasan).Where("id = ?", id).First(&expencetrasan)
	IndexRepo.DbClose(GormDB)
	
	return &expencetrasan, nil
}

func (expencetrasanRepo expencetrasanrepo) GetAll(expencetrasans []model.Expencetrasan,search *support.Search) ([]model.Expencetrasan, *httperors.HttpError) {
	results, err1 := expencetrasanRepo.Search(search, expencetrasans)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (expencetrasanRepo expencetrasanrepo) Update(id int, expencetrasan *model.Expencetrasan) (*model.Expencetrasan, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("expencetrasan with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	aexpencetrasan := model.Expencetrasan{}
	
	GormDB.Model(&aexpencetrasan).Where("id = ?", id).First(&aexpencetrasan)
	if expencetrasan.Name  == "" {
		expencetrasan.Name = aexpencetrasan.Name
	}
	if expencetrasan.Amount  == 0 {
		expencetrasan.Amount = aexpencetrasan.Amount
	}
	GormDB.Save(&expencetrasan)
	IndexRepo.DbClose(GormDB)

	return expencetrasan, nil
}
func (expencetrasanRepo expencetrasanrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := expencetrasanRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&expencetrasan).Where("id = ?", id).First(&expencetrasan)
	GormDB.Delete(expencetrasan)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (expencetrasanRepo expencetrasanrepo)ProductUserExistByid(id int) bool {
	expencetrasan := model.Expencetrasan{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&expencetrasan, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (expencetrasanRepo expencetrasanrepo) Search(Ser *support.Search, expencetrasans []model.Expencetrasan)([]model.Expencetrasan, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	expencetrasan := model.Expencetrasan{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&expencetrasan).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);
		
	// break;
	case "like":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&expencetrasans);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return expencetrasans, nil
}