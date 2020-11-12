package repository

import (
	"fmt"
	"strings"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/vcraescu/go-paginator" 
	"github.com/vcraescu/go-paginator/adapter"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)

var (
	Assetrepo assetrepo = assetrepo{}
)

///curtesy to gorm
type assetrepo struct{}

func (assetRepo assetrepo) Create(asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	if err := asset.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&asset)
	IndexRepo.DbClose(GormDB)
	return asset, nil
}
func (assetRepo assetrepo) GetOne(id int) (*model.Asset, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asset with that id does not exists!")
	}
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("Asstranss").Model(&asset).Where("id = ?", id).First(&asset)
	IndexRepo.DbClose(GormDB)
	
	return &asset, nil
}

func (assetRepo assetrepo) GetAll(assets []model.Asset,search *support.Search) ([]model.Asset, *httperors.HttpError) {
	results, err1 := assetRepo.Search(search, assets)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (assetRepo assetrepo) Update(id int, asset *model.Asset) (*model.Asset, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("asset with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	// asset := model.Asset{}
	aasset := model.Asset{}
	
	GormDB.Model(&asset).Where("id = ?", id).First(&aasset)
	if asset.Name  == "" {
		asset.Name = aasset.Name
	}
	if asset.Liscence  == "" {
		asset.Liscence = aasset.Liscence
	}
	if asset.Ownership  == "" {
		asset.Ownership = aasset.Ownership
	}
	if asset.Depreciationtype  == "" {
		asset.Depreciationtype = aasset.Depreciationtype
	}
	if asset.Depreciationrate  < 0 {
		asset.Depreciationrate = aasset.Depreciationrate
	}
	if asset.Price  < 0 {
		asset.Price = aasset.Price
	}
	if asset.ExpectedUsage  < 0 {
		asset.ExpectedUsage = aasset.ExpectedUsage
	}
	if asset.Description  == "" {
		asset.Description = aasset.Description
	}
	GormDB.Model(&asset).Where("id = ?", id).First(&asset).Update(&aasset)
	
	IndexRepo.DbClose(GormDB)

	return asset, nil
}
func (assetRepo assetrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := assetRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&asset).Where("id = ?", id).First(&asset)
	GormDB.Delete(asset)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (assetRepo assetrepo)ProductUserExistByid(id int) bool {
	asset := model.Asset{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&asset, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (assetRepo assetrepo) Search(Ser *support.Search, assets []model.Asset)([]model.Asset, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	asset := model.Asset{}
	switch(Ser.Search_operator){
	case "all":
		q := GormDB.Model(&asset).Order(Ser.Column+" "+Ser.Direction).Find(&assets)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&assets);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&assets);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Preload("Asstranss").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&assets);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
	case "like":
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&assets);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		q := GormDB.Preload("Asstranss").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&assets);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(1)
		
		if err3 := p.Results(&assets); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return assets, nil
}