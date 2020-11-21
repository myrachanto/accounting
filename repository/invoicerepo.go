package repository

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/support"
)
//Invoicerepo..
var (
	Invoicerepo invoicerepo = invoicerepo{}
)

///curtesy to gorm
type invoicerepo struct{}

func (invoiceRepo invoicerepo) Create(invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	code := invoice.Code
	t,r := Cartrepo.SumTotal(code);if r != nil {
		return nil, r
	}
	invoice.Discount = t.Discount
	invoice.Sub_total = t.Subtotal
	invoice.Total = t.Total
	tr,e := Cartrepo.CarttoTransaction(code);if e != nil {
		return nil, e
	}
	invoice.Transactions = tr 
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Create(&invoice)
	Cartrepo.DeleteAll(code)
	IndexRepo.DbClose(GormDB)
	return invoice, nil
}
func (invoiceRepo invoicerepo) View() (*model.Cinvoiceoptions, *httperors.HttpError) {
	CIOptions := &model.Cinvoiceoptions{}

	customers,err1 := Customerrepo.All()
	if err1 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	taxs,err2 := Taxrepo.All()
	if err2 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}

	products,err3 := Productrepo.All()
	if err3 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	prices,err5 := Pricerepo.All()
	if err5 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}	
	discounts,err6 := Discountrepo.All()
	if err6 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	paymentforms,err7 := Paymentformrepo.All()
	if err7 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	code,err4 := invoiceRepo.GeneCode()
	if err4 != nil {
		return nil, httperors.NewNotFoundError("You got an error fetching customers")
	}
	CIOptions.Code = code
	CIOptions.Customers = customers
	CIOptions.Taxs = taxs
	CIOptions.Products = products
	CIOptions.Prices = prices
	CIOptions.Discounts = discounts
	CIOptions.Paymentform = paymentforms
	return CIOptions, nil
} 
func (invoiceRepo invoicerepo) GetOne(id int) (*model.Invoice, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice)
	IndexRepo.DbClose(GormDB)
	
	return &invoice, nil
}
func (invoiceRepo invoicerepo) All() (t []model.Invoice, r *httperors.HttpError) {

	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Order("name").Find(&t)
	IndexRepo.DbClose(GormDB)
	return t, nil

}
func (invoiceRepo invoicerepo) GetAll(invoices []model.Invoice,search *support.Search) ([]model.Invoice, *httperors.HttpError) {
	
	results, err1 := invoiceRepo.Search(search, invoices)
	if err1 != nil {
			return nil, err1
		}
	return results, nil
}

func (invoiceRepo invoicerepo) Update(id int, invoice *model.Invoice) (*model.Invoice, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	ainvoice := model.Invoice{}
	
	GormDB.Model(&ainvoice).Where("id = ?", id).First(&ainvoice)
	// if invoice.invoice  == "" {
	// 	invoice.invoice = ainvoice.invoice
	// }
	// if invoice.Description  == "" {
	// 	invoice.Description = ainvoice.Description
	// }
	// if invoice.Subtotal  == 0 {
	// 	invoice.Subtotal = ainvoice.Subtotal
	// }
	// if invoice.Discount  == 0 {
	// 	invoice.Discount = ainvoice.Discount
	// }	
	// if invoice.AmountPaid  == 0 {
	// 	invoice.AmountPaid = ainvoice.AmountPaid
	// }
	GormDB.Save(&invoice)
	
	IndexRepo.DbClose(GormDB)

	return invoice, nil
}
func (invoiceRepo invoicerepo) Delete(id int) (*httperors.HttpSuccess, *httperors.HttpError) {
	ok := invoiceRepo.invoiceUserExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("invoice with that id does not exists!")
	}
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	GormDB.Model(&invoice).Where("id = ?", id).First(&invoice)
	GormDB.Delete(invoice)
	IndexRepo.DbClose(GormDB)
	return httperors.NewSuccessMessage("deleted successfully"), nil
}
func (invoiceRepo invoicerepo)invoiceUserExistByid(id int) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&invoice, "id =?", id)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB)
	return true
	
}
func (invoiceRepo invoicerepo)InvoiceExistByCode(code string) bool {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}

	res := GormDB.First(&invoice, "code =?", code)
	if res.Error != nil {
	   return false
	}
	IndexRepo.DbClose(GormDB) 
	return true
	
}
func (invoiceRepo invoicerepo)GeneCode() (string, *httperors.HttpError) {
	invoice := model.Invoice{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	res := GormDB.Last(&invoice)
	if res.Error != nil {
		var c1 uint = 1
		code := "CustInvNo"+strconv.FormatUint(uint64(c1), 10)
		return code, nil
	 }
	GormDB.Last(&invoice)
	c1 := invoice.ID + 1
	code := "CustInvNo"+strconv.FormatUint(uint64(c1), 10)
	IndexRepo.DbClose(GormDB)
	return code, nil
	
}
func (invoiceRepo invoicerepo) Search(Ser *support.Search, invoices []model.Invoice)([]model.Invoice, *httperors.HttpError){
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	invoice := model.Invoice{}
	switch(Ser.Search_operator){
	case "all":
	GormDB.Model(&invoice).Order(Ser.Column+" "+Ser.Direction).Find(&invoices)
		///////////////////////////////////////////////////////////////////////////////////////////////////////
		///////////////find some other paginator more effective one///////////////////////////////////////////
	
	break;
	case "equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	break;
	case "not_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "less_than" :
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "greater_than":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "less_than_or_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
	case "greater_than_ro_equal_to":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", Ser.Search_query_1).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);	
	
	break;
		 case "in":
			// db.Where("name IN (?)", []string{"myrachanto", "anto"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
		fmt.Println(s)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"(?)", s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
		break;
	 case "not_in":
			//db.Not("name", []string{"jinzhu", "jinzhu 2"}).Find(&users)
		s := strings.Split(Ser.Search_query_1,",")
	GormDB.Not(Ser.Search_column, s).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	// break;
	case "like":
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"?", "%"+Ser.Search_query_1+"%").Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	break;
	case "between":
		//db.Where("name BETWEEN ? AND ?", "lastWeek, today").Find(&users)
	GormDB.Where(Ser.Search_column+" "+Operator[Ser.Search_operator]+"? AND ?", Ser.Search_query_1, Ser.Search_query_2).Order(Ser.Column+" "+Ser.Direction).Find(&invoices);
	
	   break;
	default:
	return nil, httperors.NewNotFoundError("check your operator!")
	}
	IndexRepo.DbClose(GormDB)
	
	return invoices, nil
}