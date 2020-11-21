package repository

import (
	"fmt"
	"strings"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Sinvoicerepo ...
var (
	Sinvoicerepo sinvoicerepo = sinvoicerepo{}
)

///curtesy to gorm
type sinvoicerepo struct{}

func (sinvoiceRepo sinvoicerepo) Create(sinvoice *model.SInvoice) (*httperors.HttpSuccess, *httperors.HttpError) {
	code := sinvoice.Code
	t,r := Cartrepo.SumTotal(code);if r != nil {
		return nil, r
	}
	sinvoice.Discount = t.Discount
	sinvoice.Sub_total = t.Subtotal
	sinvoice.Total = t.Total 
	tr,e := Scartrepo.CarttoTransaction(code);if e != nil {
		return nil, e
	}
	//invoice.Transactions = tr 
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&sinvoice)
	for _, transact := range tr{
		GormDB.Model(model.SInvoice{}).Association("STransactions").Append(transact)
	}
	s,f := Scartrepo.DeleteAll(code)
	if f != nil {
		return nil,f
	}
	IndexRepo.DbClose(GormDB)
	return s, nil
}
func (sinvoiceRepo sinvoicerepo) GetOne(id int) (*model.SInvoice, *httperors.HttpError) {
	ok := sinvoiceRepo.sinvoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("sinvoice with that id does not exists!")
	}
	sinvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Preload("STransactions").Model(&sinvoice).Where("id = ?", id).First(&sinvoice)
	IndexRepo.DbClose(GormDB)
	
	return &sinvoice, nil
}

func (sinvoiceRepo sinvoicerepo) GetAll(sinvoices []model.SInvoice,search *support.Search) ([]model.SInvoice, *httperors.HttpError) {
	
	results, err1 := sinvoiceRepo.Search(search, sinvoices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}
func (sinvoiceRepo sinvoicerepo) View() (*model.Soptions, *httperors.HttpError) {
	code,err := sinvoiceRepo.GeneCode();if err != nil {
		return nil, err
	} 
	c,e := Supplierrepo.GetOptions();if e != nil {
		return nil,e
	}
	m,me := Productrepo.GetOptions();if me != nil {
		return nil,me
	}
	options := model.Soptions{}
	options.Code = code
	options.Supplier = c
	options.Product = m
	return &options, nil
}
// func (sinvoiceRepo sinvoicerepo) Update(id int, sinvoice *model.SInvoice) (*model.SInvoice, *httperors.HttpError) {
// 	ok := sinvoiceRepo.sinvoiceUserExistByid(id)
// 	if !ok {
// 		return nil, httperors.NewNotFoundError("sinvoice with that id does not exists!")
// 	}
	
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return nil, err1
// 	}
// 	asinvoice := model.SInvoice{}
	
// 	GormDB.Model(&asinvoice).Where("id = ?", id).First(&asinvoice)
// 	// if sinvoice.sinvoice  == "" {
// 	// 	sinvoice.sinvoice = asinvoice.sinvoice
// 	// }
// 	// if sinvoice.Description  == "" {
// 	// 	sinvoice.Description = asinvoice.Description
// 	// }
// 	// if sinvoice.Subtotal  == 0 {
// 	// 	sinvoice.Subtotal = asinvoice.Subtotal
// 	// }
// 	// if sinvoice.Discount  == 0 {
// 	// 	sinvoice.Discount = asinvoice.Discount
// 	// }	
// 	// if sinvoice.AmountPaid  == 0 {
// 	// 	sinvoice.AmountPaid = asinvoice.AmountPaid
// 	// }
// 	GormDB.Model(&sinvoice).Where("id = ?", id).First(&sinvoice).Update(&asinvoice)
	
// 	IndexRepo.DbClose(GormDB)

// 	return sinvoice, nil
// }
// func (sinvoiceRepo sinvoicerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
// 	ok := sinvoiceRepo.sinvoiceUserExistByid(id)
// 	if !ok {
// 		return nil, httperors.NewNotFoundError("sinvoice with that id does not exists!")
// 	}
// 	sinvoice := model.SInvoice{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return nil, err1
// 	}
// 	GormDB.Model(&sinvoice).Where("id = ?", id).First(&sinvoice)
// 	GormDB.Delete(sinvoice)
// 	IndexRepo.DbClose(GormDB)
// 	return httperors.NewSuccessMessage("deleted successfully"), nil
// }
func (sinvoiceRepo sinvoicerepo)sinvoiceUserExistByid(id int) bool {
	sinvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	} 
	res := GormDB.First(&sinvoice, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (sinvoiceRepo sinvoicerepo)InvoiceExistByCode(code string) bool {
	sinvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&sinvoice, "code =?", code)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB) 
	return true
	
}
func (sinvoiceRepo sinvoicerepo)GeneCode() (string, *httperors.HttpError) {
	sinvoice := model.SInvoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	res := GormDB.Last(&sinvoice)
	if res.Error != nil {
		c1 := 1
		code := "SuppInvNo_"+ string(c1)
		return code, nil
	 }
	GormDB.Last(&sinvoice)
	c1 := sinvoice.ID + 1
	code := "SuppInvNo_"+ string(c1)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (sinvoiceRepo sinvoicerepo) Search(Ser *support.Search, sinvoices []model.SInvoice)([]model.SInvoice, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	sinvoice := model.SInvoice{}
	switch(Ser.Search_operator){
	case "all":
		GormDB.Preload("STransactions").Model(&sinvoice).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
		
	break;
	case "equal_to":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);
		
	break;
	case "not_equal_to":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);	
		
	break;
	case "less_than" :
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);	
		
	break;
	case "greater_than":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);	
		
	break;
	case "less_than_or_equal_to":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);	
		
	break;
	case "greater_than_ro_equal_to":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);	
		
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);
		
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		GormDB.Preload("STransactions").Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);
		
	// break;
	case "like":
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);
		
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
		GormDB.Preload("STransactions").Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&sinvoices);
		
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return sinvoices, nil
}