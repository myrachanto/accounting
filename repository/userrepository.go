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
//Userrepo ...
var (
	Userrepo userrepo = userrepo{}
)

///curtesy to gorm
type userrepo struct{}

func (userRepo userrepo) Create(user *model.User) (string, *httperors.HttpError) {
	if err := user.Validate(); err != nil {
		return "", err
	}
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return "", err1
	}
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return "", httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = userRepo.UserExist(user.Email)
	if ok {
		return "", httperors.NewNotFoundError("Your email already exists!")
	}
	
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return "", err2
	}
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	
	fmt.Println(user)
	GormDB.Create(&user)
	IndexRepo.DbClose(GormDB)
	return "user created successifully", nil
}
func (userRepo userrepo) Login(auser *model.LoginUser) (*model.Auth, *httperors.HttpError) {
	if err := auser.Validate(); err != nil {
		return nil, err
	}
	ok := userRepo.UserExist(auser.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	user := model.User{}

	GormDB.Model(&user).Where("email = ?", auser.Email).First(&user)
	ok = user.Compare(auser.Password, user.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.Token{
		UserID: user.ID,
		UName: user.UName,
		Admin:user.Admin,
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
	// messages ,e := userRepo.UnreadMessages(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := userRepo.UnreadNortifications(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.Auth{UserID:user.ID, UName:user.UName,Admin:user.Admin, Picture:user.Picture, Token:tokenString}
	GormDB.Create(&auth)
	IndexRepo.DbClose(GormDB)
	
	return auth, nil
}
func (userRepo userrepo) Logout(token string) (*httperors.HttpError) {
	auth := model.Auth{}
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
func (userRepo userrepo) All() (t []model.User, r *httperors.HttpError) {

	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&user).Order("name").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (userRepo userrepo) GetOne(id int) (*model.User, *httperors.HttpError) {
	ok := userRepo.UserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&user).Where("id = ?", id).First(&user)
	IndexRepo.DbClose(GormDB)
	
	return &user, nil
}

func (userRepo userrepo) GetAll(users []model.User,search *support.Search) ([]model.User, *httperors.HttpError) {
	results, err1 := userRepo.Search(search, users)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (userRepo userrepo) Update(id int, user *model.User) (*model.User, *httperors.HttpError) {
	ok := userRepo.UserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}
	
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	User := model.User{}
	uuser := model.User{}
	
	GormDB.Model(&User).Where("id = ?", id).First(&uuser)
	if user.FName  == "" {
		user.FName = uuser.FName
	}
	if user.LName  == "" {
		user.LName = uuser.LName
	}
	if user.UName  == "" {
		user.UName = uuser.UName
	}
	if user.Phone  == "" {
		user.Phone = uuser.Phone
	}
	if user.Address  == "" {
		user.Address = uuser.Address
	}
	if user.Picture  == "" {
		user.Picture = uuser.Picture
	}
	if user.Email  == "" {
		user.Email = uuser.Email
	}
	if user.Admin  == false {
		user.Admin = true
	}
	GormDB.Save(&user)
	
	IndexRepo.DbClose(GormDB)

	return user, nil
}
func (userRepo userrepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := userRepo.UserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&user).Where("id = ?", id).First(&user)
	GormDB.Delete(user)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (userRepo userrepo)UserExist(email string) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&user, "email =?", email)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (userRepo userrepo)UserExistByid(id int) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&user, "id =?", id)
	if res.Error != nil{
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
// func (userRepo userrepo)UnreadMessages(id uint)  (int, *httperors.HttpError)  {
// 	messages := []model.Message{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&messages)	
// 	 c := 0
// 	 for i, _:= range messages{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil
	
// }
// func (userRepo userrepo)UnreadNortifications(id uint)  (int, *httperors.HttpError)  {
// 	ns := []model.Nortification{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&ns)	
// 	 c := 0
// 	 for i, _:= range ns{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil
	
// }


func (userRepo userrepo) Search(Ser *support.Search, users []model.User)([]model.User, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	user := model.User{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Model(&user).Order(Ser.Column+" "+Ser.Direction).Find(&users)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);
		
	break;
	case "not_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);	
		
	break;
	case "less_than" :
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);	
		
	break;
	case "greater_than":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&users);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Order(Ser.Column+" "+Ser.Direction).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&users);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Order(Ser.Column+" "+Ser.Direction).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&users);
		
	// break;
case "like":
	// fmt.Println(Ser.Search_query_1)
	if Ser.Search_query_1 == "all" {
			//db.Order("name DESC")
	GormDB.Order(Ser.Column+" "+Ser.Direction).Find(&users)

	}else {

		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&users);
	
	}
break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Order(Ser.Column+" "+Ser.Direction).Find(&users)
		GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&users);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return users, nil
}