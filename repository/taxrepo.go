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
	Taxrepo taxrepo = taxrepo{}
)

///curtesy to gorm
type taxrepo struct{}

func (taxRepo taxrepo) Create(tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	if err := tax.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&tax)
	IndexRepo.DbClose(GormDB)
	return tax, nil
}
func (taxRepo taxrepo) GetOne(id int) (*model.Tax, *httperors.HttpError) {
	ok := taxRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("tax with that id does not exists!")
	}
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&tax).Where("id = ?", id).First(&tax)
	IndexRepo.DbClose(GormDB)
	
	return &tax, nil
}
func (taxRepo taxrepo) All() (t []model.Tax, r *httperors.HttpError) {

	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&tax).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (taxRepo taxrepo) GetOption()([]model.Tax, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	tax := model.Tax{}
	taxs := []model.Tax{}
	GormDB.Model(&tax).Find(&taxs)
	return taxs, nil
}
func (taxRepo taxrepo) GetAll(taxs []model.Tax,search *support.Search) ([]model.Tax, *httperors.HttpError) {

	results, err1 := taxRepo.Search(search, taxs)
	if err1 != nil {
			return nil, err1
		}
		fmt.Println(results)
	return results, nil
}

func (taxRepo taxrepo) Update(id int, tax *model.Tax) (*model.Tax, *httperors.HttpError) {
	ok := taxRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("tax with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	atax := model.Tax{}
	
	GormDB.Model(&atax).Where("id = ?", id).First(&atax)
	if tax.Name  == "" {
		tax.Name = atax.Name
	}
	if tax.Title  == "" {
		tax.Title = atax.Title
	}
	if tax.Description  == "" {
		tax.Description = atax.Description
	}
	fmt.Println(tax)
	GormDB.Model(&atax).Where("id = ?", id).Update(&tax)
	
	IndexRepo.DbClose(GormDB)

	return tax, nil
}
func (taxRepo taxrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := taxRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&tax).Where("id = ?", id).First(&tax)
	GormDB.Delete(tax)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (taxRepo taxrepo)ProductUserExistByid(id int) bool {
	tax := model.Tax{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&tax, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (taxRepo taxrepo) Search(Ser *support.Search, taxs []model.Tax)([]model.Tax, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	tax := model.Tax{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&tax).Order(Ser.Column+" "+Ser.Direction).Find(&taxs)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&taxs)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		fmt.Println(p.Results(&taxs))
				if err3 := p.Results(&taxs); err3 != nil {
					return nil, httperors.NewNotFoundError("something went wrong paginating!")
				}

		}else {

			q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&taxs);
			p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
			p.SetPage(Ser.Page)
			fmt.Println(p.Results(&taxs))
			if err3 := p.Results(&taxs); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&taxs);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&taxs); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return taxs, nil
}