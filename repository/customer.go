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
)
//Customerrepo ...
var (
	Customerrepo customerrepo = customerrepo{}
)

///curtesy to gorm
type customerrepo struct{}

func (customerRepo customerrepo) Create(customer *model.Customer) (string, *httperors.HttpError) {
	if err := customer.Validate(); err != nil {
		return "", err
	}
	ok, err1 := customer.ValidatePassword(customer.Password)
	if !ok {
		return "", err1
	}
	ok = customer.ValidateEmail(customer.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = customerRepo.customerExist(customer.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	hashpassword, err2 := customer.HashPassword(customer.Password)
	if err2 != nil {
		return "", err2
	}
	customer.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	
	fmt.Println(customer)
	GormDB.Create(&customer)
	IndexRepo.DbClose(GormDB)
	return "customer created successifully", nil
}
func (customerRepo customerrepo) Login(acustomer *model.Logincustomer) (*model.CustomnerAuth, *httperors.HttpError) {
	if err := acustomer.Validate(); err != nil {
		return nil, err
	}
	ok := customerRepo.customerExist(acustomer.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	GormDB.Model(&customer).Where("email = ?", acustomer.Email).First(&customer)
	ok = customer.Compare(acustomer.Password, customer.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.CustomerToken{
		CustomerID: customer.ID,
		Name: customer.Name,
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
	// messages ,e := customerRepo.UnreadMessages(customer.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := customerRepo.UnreadNortifications(customer.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.CustomnerAuth{CustomerID:customer.ID, Name:customer.Name, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}
func (customerRepo customerrepo) Logout(token string) (*httperors.HttpError) {
	auth := model.CustomnerAuth{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	res := GormDB.First(&auth, "token =?", token)
	if res.Error != nil {
		return httperors.NewNotFoundError("Something went wrong logging out!")
	 }
	
	GormDB.Model(&auth).Where("token =?", token).First(&auth)
	
	GormDB.Delete(auth)
	IndexRepo.DbClose(GormDB)
	
	return  nil
}
func (customerRepo customerrepo) Forgot(email string) (string, *httperors.HttpError) {
	ok := customerRepo.customerExist(email)
	if !ok {
		return "", httperors.NewNotFoundError("That Email does not exists with our records!")
	}
	
	return "Email sent!", nil
}
func (customerRepo customerrepo) GetOne(id int) (*model.Customer, *httperors.HttpError) {
	ok := customerRepo.customercustomerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	fmt.Println(customer)
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	IndexRepo.DbClose(GormDB)
	
	return &customer, nil
}
func (customerRepo customerrepo) GetOptions()([]model.Customer, *httperors.HttpError){

	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	customers := []model.Customer{}
	GormDB.Model(&customer).Find(&customers)
	return customers, nil
}
func (customerRepo customerrepo) GetAll(search *support.Search) ([]model.Customer,*httperors.HttpError) {
	customers := []model.Customer{} 
	results, err1 := customerRepo.Search(search, customers)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (customerRepo customerrepo) All() (t []model.Customer, r *httperors.HttpError) {

	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&customer).Order("name").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}

// func (customerRepo customerrepo) GetAll(search *support.Search) ([]interface{}, *httperors.HttpError) {
// 	customer := model.Customer{}
// 	// customers := []model.Customer{}
// 	// results, err1 := customerRepo.Search(search, customer)
// 	 results, err1 := support.SearchQuery(search, customer)
// 	if err1 != nil {
// 			return nil, err1
// 		}
// 	return results, nil 
// }

func (customerRepo customerrepo) Update(id int, customer *model.Customer) (*model.Customer, *httperors.HttpError) {
	ok := customerRepo.customercustomerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	acustomer := model.Customer{}
	
	GormDB.Model(&acustomer).Where("id = ?", id).First(&acustomer)
	if customer.Name  == "" {
		customer.Name = acustomer.Name
	}
	if customer.Company  == "" {
		customer.Company = acustomer.Company
	}
	if customer.Phone  == "" {
		customer.Phone = acustomer.Phone
	}
	if customer.Email  == "" {
		customer.Email = acustomer.Email
	}
	if customer.Address  == "" {
		customer.Address = acustomer.Address
	}
	if customer.Picture  == "" {
		customer.Picture = acustomer.Picture
	}
	GormDB.Save(&customer)
	
	IndexRepo.DbClose(GormDB)

	return customer, nil
}
func (customerRepo customerrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := customerRepo.customercustomerExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("customer with that id does not exists!")
	}
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	} 
	GormDB.Model(&customer).Where("id = ?", id).First(&customer)
	GormDB.Delete(customer)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (customerRepo customerrepo)customercustomerExistByid(id int) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&customer, "id =?", id)
		if res.Error != nil {
			return false
		}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo)customerExist(email string) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&customer, "email =?", email)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo)customerExistByid(id int) bool {
	customer := model.Customer{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&customer, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (customerRepo customerrepo) Search(Ser *support.Search, customers []model.Customer)([]model.Customer, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	customer := model.Customer{}
	// // invoices := model.Invoice{}
	// fmt.Println(&customer)
	switch(Ser.Search_operator){
	case "all":
		//db.Order("name DESC")
		GormDB.Model(&customer).Order(Ser.Column+" "+Ser.Direction).Find(&customers)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		
	break;
	case "not_equal_to":
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		
	break;
	case "less_than" :
		// order := &Order
		// db.Where("id = ? and status = ?", reqOrder.id, "cart")
		// .Preload("OrderItems").Preload("OrderItems.Item").First(&order)
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		
	break;
	case "greater_than":
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&customers);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&customers)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&customers)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Preload("Invoices").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&customers)

	}else {

		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&customers);
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&customers)
		GormDB.Preload("Invoices").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&customers);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return customers, nil
}
////////////subject to futher scrutiny/////////////////////////////////
// func (customerRepo customerrepo)paginator(q *gorm.DB, Ser *support.Search, customers []model.Customer) ([]model.Customer, *httperors.HttpError) {
// 	p := paginator.New(adapter.NewGORMAdapter(q), Ser.Per_page)
// 	p.SetPage(Ser.Page)
// 	// fmt.Println(Ser.Per_page)
// 	err3 := p.Results(&customers)
// 	if err3 != nil {
// 		return nil, httperors.NewNotFoundError("something went wrong paginating!")
// 	}
// 	return customers, nil
	
// }