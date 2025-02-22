package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Productrepo ...
var (
	Productrepo productrepo = productrepo{}
)

///curtesy to gorm
type productrepo struct{}

func (productRepo productrepo) Create(product *model.Product) (*model.Product, *httperors.HttpError) {
	if err := product.Validate(); err != nil {
		return nil, err
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cat, err := Categoryrepo.GetMajorcat(product.Category)
	if err != nil {
		return nil, httperors.NewNotFoundError("category with that name does not exists!")
	}
	product.Majorcategory = cat.Majorcategory 
	GormDB.Create(&product)
	IndexRepo.DbClose(GormDB)
	return product, nil
}
func (productRepo productrepo) GetOne(id int) (*model.Product, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	IndexRepo.DbClose(GormDB)
	
	return &product, nil
}

func (productRepo productrepo) View() ([]model.Category, *httperors.HttpError) {
	mc, e := Categoryrepo.All()
	if e != nil{
		return nil, e
	}
	return mc, nil
}

func (productRepo productrepo) All() (t []model.Product, r *httperors.HttpError) {

	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&product).Order("name").Find(&t)
	
	IndexRepo.DbClose(GormDB)
	return t, nil

}

func (productRepo productrepo) GetProducts(products []model.Product,search *support.Productsearch) ([]model.Product, *httperors.HttpError) {
	results, err1 := productRepo.SearchFront(search, products)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (productRepo productrepo) GetAll(products []model.Product,search *support.Search) ([]model.Product, *httperors.HttpError) {
	results, err1 := productRepo.Search(search, products)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (productRepo productrepo) Update(id int, product *model.Product) (*model.Product, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("product with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	cat, err := Categoryrepo.GetMajorcat(product.Category)
	if err != nil {
		return nil, httperors.NewNotFoundError("category with that name does not exists!")
	}
	product.Majorcategory = cat.Majorcategory 
	aproduct := model.Product{}
	
	GormDB.Model(&aproduct).Where("id = ?", id).First(&aproduct)
	if product.Name  == "" {
		product.Name = aproduct.Name
	}
	if product.Title  == "" {
		product.Title = aproduct.Title
	}
	if product.Description  == "" {
		product.Description = aproduct.Description
	}
	if product.Category  == "" {
		product.Category = aproduct.Category
	}
	if product.Majorcategory  == "" { 
		product.Majorcategory = aproduct.Majorcategory
	}
	if product.Quantity  == 0 { 
		product.Quantity = aproduct.Quantity
	}
	GormDB.Save(&product)
	
	
	IndexRepo.DbClose(GormDB)

	return product, nil
}
func (productRepo productrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := productRepo.ProductUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("Product with that id does not exists!")
	}
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&product).Where("id = ?", id).First(&product)
	GormDB.Delete(product)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (productRepo productrepo)ProductUserExistByid(id int) bool {
	product := model.Product{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&product, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}

func (productRepo productrepo) GetOptions()([]model.Product, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	products := []model.Product{}
	GormDB.Model(&product).Find(&products)
	return products, nil
}
func (productRepo productrepo) Search(Ser *support.Search, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	// break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	

		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&products);
			
		}
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}

func (productRepo productrepo) SearchFront(Ser *support.Productsearch, products []model.Product)([]model.Product, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	product := model.Product{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&product).Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		///product = %as% AND price (>/=/<=/>=/between)
		//db.Where("name = ? AND age >= ? ", "myrachanto", "28").Find(&users)
		//db.Where("name LIKE ?", "%a%").Find(&users)
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than" :
		fmt.Println(Ser)
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Column+" LIKE ? AND " +Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Name+"%", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&products);	
		
	break;
	case "like":
		// fmt.Println(Ser.Search_query_1)
		if Ser.Search_query_1 == "all" {
				//db.Order("name DESC")
		GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&products)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		

		}else {

			GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&products);
		
		}
	break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return products, nil
}