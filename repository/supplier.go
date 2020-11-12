package repository

import (
	"fmt"
	"log"
	"os"
	"strings"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
	"github.com/vcraescu/go-paginator"
	"github.com/vcraescu/go-paginator/adapter"
)
  
var (
	Supplierrepo supplierrepo = supplierrepo{}
)

///curtesy to gorm
type supplierrepo struct{}

func (supplierRepo supplierrepo) Create(supplier *model.Supplier) (string, *httperors.HttpError) {
	if err := supplier.Validate(); err != nil {
		return "", err
	}
	ok, err1 := supplier.ValidatePassword(supplier.Password)
	if !ok {
		return "", err1
	}
	ok = supplier.ValidateEmail(supplier.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = supplierRepo.supplierExist(supplier.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	hashpassword, err2 := supplier.HashPassword(supplier.Password)
	if err2 != nil {
		return "", err2
	}
	supplier.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	
	fmt.Println(supplier)
	GormDB.Create(&supplier)
	IndexRepo.DbClose(GormDB)
	return "supplier created successifully", nil
}
func (supplierRepo supplierrepo) Login(asupplier *model.Loginsupplier) (*model.SupplierAuth, *httperors.HttpError) {
	if err := asupplier.Validate(); err != nil {
		return nil, err
	}
	ok := supplierRepo.supplierExist(asupplier.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	supplier := model.Supplier{}
	GormDB.Model(&supplier).Where("email = ?", asupplier.Email).First(&supplier)
	ok = supplier.Compare(asupplier.Password, supplier.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.SupplierToken{
		SupplierID: supplier.ID,
		Name: supplier.Name,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: model.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading key")
	}
	encyKey := os.Getenv("EncryptionKey")
	tokenString, error := token.SignedString([]byte(encyKey))
	if error != nil {
		fmt.Println(error)
	}
	// messages ,e := supplierRepo.UnreadMessages(supplier.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := supplierRepo.UnreadNortifications(supplier.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.SupplierAuth{SupplierID:supplier.ID, Name:supplier.Name, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}
func (supplierRepo supplierrepo) Logout(token string) (*httperors.HttpError) {
	auth := model.CustomnerAuth{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	if GormDB.First(&auth, "token =?", token).RecordNotFound(){
		return httperors.NewNotFoundError("Something went wrong logging out!")
	 }
	
	GormDB.Model(&auth).Where("token =?", token).First(&auth)
	
	GormDB.Delete(auth)
	IndexRepo.DbClose(GormDB)
	
	return  nil
}
func (supplierRepo supplierrepo) GetOptions()([]model.Supplier, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	supplier := model.Supplier{}
	suppliers := []model.Supplier{}
	GormDB.Model(&supplier).Find(&suppliers)
	return suppliers, nil
}
func (supplierRepo supplierrepo) Forgot(email string) (string, *httperors.HttpError) {
	ok := supplierRepo.supplierExist(email)
	if !ok {
		return "", httperors.NewNotFoundError("That Email does not exists with our records!")
	}
	
	return "Email sent!", nil
}
func (supplierRepo supplierrepo) GetOne(id int) (*model.Supplier, *httperors.HttpError) {
	ok := supplierRepo.suppliersupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(supplier)
	GormDB.Model(&supplier).Where("id = ?", id).First(&supplier)
	IndexRepo.DbClose(GormDB)
	
	return &supplier, nil
}
func (supplierRepo supplierrepo) GetAll(search *support.Search) ([]model.Supplier,*httperors.HttpError) {
	suppliers := []model.Supplier{} 
	results, err1 := supplierRepo.Search(search, suppliers)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (supplierRepo supplierrepo) All() (t []model.Supplier, r *httperors.HttpError) {

	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	q := GormDB.Model(&supplier).Order("name").Find(&t)
	p := paginator.New(adapter.NewGORMAdapter(q), 40)
	p.SetPage(1)

	
	if err3 := p.Results(&t); err3 != nil {
		return nil, httperors.NewNotFoundError("something went wrong paginating!")
	}
	IndexRepo.DbClose(GormDB)
	return t, nil

}

// func (supplierRepo supplierrepo) GetAll(search *support.Search) ([]interface{}, *httperors.HttpError) {
// 	supplier := model.Supplier{}
// 	// suppliers := []model.Supplier{}
// 	// results, err1 := supplierRepo.Search(search, supplier)
// 	 results, err1 := support.SearchQuery(search, supplier)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil 
// }

func (supplierRepo supplierrepo) Update(id int, supplier *model.Supplier) (*model.Supplier, *httperors.HttpError) {
	ok := supplierRepo.suppliersupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	asupplier := model.Supplier{}
	
	GormDB.Model(&asupplier).Where("id = ?", id).First(&asupplier)
	if supplier.Name  == "" {
		supplier.Name = asupplier.Name
	}
	if supplier.Company  == "" {
		supplier.Company = asupplier.Company
	}
	if supplier.Phone  == "" {
		supplier.Phone = asupplier.Phone
	}
	if supplier.Email  == "" {
		supplier.Email = asupplier.Email
	}
	if supplier.Address  == "" {
		supplier.Address = asupplier.Address
	}
	if supplier.Picture  == "" {
		supplier.Picture = asupplier.Picture
	}
	GormDB.Model(&asupplier).Where("id = ?", id).Update(&supplier)
	
	IndexRepo.DbClose(GormDB)

	return supplier, nil
}
func (supplierRepo supplierrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := supplierRepo.suppliersupplierExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("supplier with that id does not exists!")
	}
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&supplier).Where("id = ?", id).First(&supplier)
	GormDB.Delete(supplier)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (supplierRepo supplierrepo)suppliersupplierExistByid(id int) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&supplier, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (supplierRepo supplierrepo)supplierExist(email string) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&supplier, "email =?", email).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (supplierRepo supplierrepo)supplierExistByid(id int) bool {
	supplier := model.Supplier{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	if GormDB.First(&supplier, "id =?", id).RecordNotFound(){
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (supplierRepo supplierrepo) Search(Ser *support.Search, suppliers []model.Supplier)([]model.Supplier, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	supplier := model.Supplier{}
	// // invoices := model.Invoice{}
	// fmt.Println(&supplier)
	switch(Ser.Search_operator){
	case "all":
		//db.Order("name DESC")
		q := GormDB.Model(&supplier).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "not_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than" :
		// order := &Order
		// db.Where("id = ? and status = ?", reqOrder.id, "cart")
		// .Preload("OrderItems").Preload("OrderItems.Item").First(&order)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "less_than_or_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
	case "greater_than_ro_equal_to":
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);	
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&suppliers)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&suppliers)
		s := strings.Split(Ser.Search_query_1,",")
		q := GormDB.Preload("Invoices").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	q := GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&suppliers)
	///////////////////////////////////////////////////////////////////////////////////////////////////////
	///////////////find some other paginator more effective one///////////////////////////////////////////
	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
	p.SetPage(Ser.Page)
	
	fmt.Println(p.Results(&suppliers))
			if err3 := p.Results(&suppliers); err3 != nil {
				return nil, httperors.NewNotFoundError("something went wrong paginating!")
			}

	}else {

		q := GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		fmt.Println(p.Results(&suppliers,))
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&suppliers)
		q := GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&suppliers);
		p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
		p.SetPage(Ser.Page)
		
		if err3 := p.Results(&suppliers); err3 != nil {
			return nil, httperors.NewNotFoundError("something went wrong paginating!")
		}
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return suppliers, nil
}
////////////subject to futher scrutiny/////////////////////////////////
// func (supplierRepo supplierrepo)paginator(q *gorm.DB, Ser *support.Search, suppliers []model.Supplier) ([]model.Supplier, *httperors.HttpError) {
// 	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
// 	p.SetPage(Ser.Page)
// 	// fmt.Println(Ser.Per_page)
// 	err3 := p.Results(&suppliers)
// 	if err3 != nil {
// 		return nil, httperors.NewNotFoundError("something went wrong paginating!")
// 	}
// 	return suppliers, nil
	
// }