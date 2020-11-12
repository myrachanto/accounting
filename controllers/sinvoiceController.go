package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
	"github.com/myrachanto/accounting/support"
)
 
var (
	SinvoiceController sinvoiceController = sinvoiceController{}
)
type sinvoiceController struct{ }
/////////controllers/////////////////
func (controller sinvoiceController) Create(c echo.Context) error {
	sinvoice := &model.SInvoice{}
	
	if err := c.Bind(sinvoice); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	s, f := service.Sinvoiceservice.Create(sinvoice)
	if f != nil {
		return c.JSON(f.Code, f)
	}
	return c.JSON(http.StatusCreated, s)
}

func (controller sinvoiceController) View(c echo.Context) error {
	code, problem := service.Invoiceservice.View()
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, code)	
}
func (controller sinvoiceController) GetAll(c echo.Context) error {
	sinvoices := []model.SInvoice{}
	column := string(c.QueryParam("column"))
	direction := string(c.QueryParam("direction"))
	search_column := string(c.QueryParam("search_column"))
	search_operator := string(c.QueryParam("search_operator"))
	search_query_1 := string(c.QueryParam("search_query_1"))
	search_query_2 := string(c.QueryParam("search_query_2"))
	per_page, err := strconv.Atoi(c.QueryParam("per_page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid per number")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println("------------------------")
	search := &support.Search{Column:column, Direction:direction,Search_column:search_column,Search_operator:search_operator,Search_query_1:search_query_1,Search_query_2:search_query_2,Per_page:per_page}
	
	sinvoices, err3 := service.Sinvoiceservice.GetAll(sinvoices,search)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, sinvoices)
} 
func (controller sinvoiceController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	sinvoice, problem := service.Sinvoiceservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, sinvoice)	
}

// func (controller sinvoiceController) Update(c echo.Context) error {
		
// 	sinvoice :=  &model.SInvoice{}
// 	if err := c.Bind(sinvoice); err != nil {
// 		httperror := httperors.NewBadRequestError("Invalid json body")
// 		return c.JSON(httperror.Code, httperror)
// 	}	
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		httperror := httperors.NewBadRequestError("Invalid ID")
// 		return c.JSON(httperror.Code, httperror)
// 	}
// 	updatedsinvoice, problem := service.Sinvoiceservice.Update(id, sinvoice)
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, updatedsinvoice)
// }

// func (controller sinvoiceController) Delete(c echo.Context) error {
// 	id, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		httperror := httperors.NewBadRequestError("Invalid ID")
// 		return c.JSON(httperror.Code, httperror)
// 	}
// 	success, failure := service.Sinvoiceservice.Delete(id)
// 	if failure != nil {
// 		return c.JSON(failure.Code, failure)
// 	}
// 	return c.JSON(success.Code, success)
		
// }
