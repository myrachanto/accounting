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
	Discountrepo discountrepo = discountrepo{}
)

///curtesy to gorm
type discountrepo struct{}

func (discountRepo discountrepo) Create(discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	if err := discount.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&discount)
	IndexRepo.DbClose(GormDB)
	return discount, nil
}
func (discountRepo discountrepo) All() (t []model.Discount, r *httperors.HttpError) {

	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&discount).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (discountRepo discountrepo) GetOne(id int) (*model.Discount, *httperors.HttpError) {
	ok := discountRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("discount with that id does not exists!")
	}
	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&discount).Where("id = ?", id).First(&discount)
	IndexRepo.DbClose(GormDB)
	
	return &discount, nil
}
func (discountRepo discountrepo) GetOption(id int)([]model.Discount, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	discounts := []model.Discount{}
	GormDB.Where("id = ? AND buy = ? ", id, true).Find(&discounts)
	return discounts, nil
}
func (discountRepo discountrepo) GetOptionsell(id int)([]model.Discount, *httperors.HttpError){
	ok := Productrepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	discounts := []model.Discount{}
	GormDB.Where("id = ? AND buy = ? ", id, false).Find(&discounts)
	return discounts, nil
}
func (discountRepo discountrepo) GetAll(discounts []model.Discount,search *support.Search) ([]model.Discount, *httperors.HttpError) {
	results, err1 := discountRepo.Search(search, discounts)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (discountRepo discountrepo) Update(id int, discount *model.Discount) (*model.Discount, *httperors.HttpError) {
	ok := discountRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("discount with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	adiscount := model.Discount{}
	
	GormDB.Model(&adiscount).Where("id = ?", id).First(&adiscount)
	if discount.Name  == "" {
		discount.Name = adiscount.Name
	}
	if discount.Title  == "" {
		discount.Title = adiscount.Title
	}
	if discount.Description  == "" {
		discount.Description = adiscount.Description
	}
	fmt.Println(discount)
	GormDB.Model(&adiscount).Where("id = ?", id).Update(&discount)

	return discount, nil
}
func (discountRepo discountrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := discountRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&discount).Where("id = ?", id).First(&discount)
	GormDB.Delete(discount)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (discountRepo discountrepo)ProductUserExistByid(id int) bool {
	discount := model.Discount{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&discount, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (discountRepo discountrepo) Search(Ser *support.Search, discounts []model.Discount)([]model.Discount, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	discount := model.Discount{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&discount).Order(Ser.Column+" "+Ser.Direction).Find(&discounts)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
case "like":
	fmt.Println(Ser)
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&discounts)
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////find some other paginator more effective one///////////////////////////////////////////
	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
	p.SetPage(Ser.Page)
	
	fmt.Println(p.Results(&discounts))
			if err3 := p.Results(&discounts); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}

	}else {

		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&discounts);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		fmt.Println(p.Results(&discounts))
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&discounts);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&discounts); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return discounts, nil
}