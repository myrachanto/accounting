package controllers

import(
	"fmt"
	"strconv"	
	"net/http"
	"github.com/labstack/echo"
	"github.com/myrachanto/accounting/httperors"
	"github.com/myrachanto/accounting/model"
	"github.com/myrachanto/accounting/service"
)
 
var (
	CartController cartController = cartController{}
)
type cartController struct{ }
/////////controllers/////////////////
func (controller cartController) Create(c echo.Context) error {
	cart := &model.Cart{}
	
	cart.Name = c.FormValue("name")
	cart.Code = c.FormValue("code")
	fmt.Println(cart.Code)
	fmt.Println(c.FormValue("quantity"), "qty")
	q, err := strconv.ParseFloat(c.FormValue("quantity"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid quantity")
		return c.JSON(httperror.Code, httperror)
	}
	f, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid price")
		return c.JSON(httperror.Code, httperror)
	}
	t, err := strconv.ParseFloat(c.FormValue("tax"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid tax")
		return c.JSON(httperror.Code, httperror)
	}

	d, err := strconv.ParseFloat(c.FormValue("discount"), 64)
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid discount")
		return c.JSON(httperror.Code, httperror)
	}
	cart.Price = f
	cart.Quantity= q
	cart.Discount = d
	cart.Tax = t
	createdcart, err1 := service.Cartservice.Create(cart)
	if err1 != nil {
		return c.JSON(err1.Code, err1)
	}
	return c.JSON(http.StatusCreated, createdcart)
}

func (controller cartController) View(c echo.Context) error {
	code := c.Param("code")
	options, problem := service.Cartservice.View(code)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, options)	
}
func (controller cartController) GetAll(c echo.Context) error {
	carts := []model.Cart{}
	carts, err3 := service.Cartservice.GetAll(carts)
	if err3 != nil {
		return c.JSON(err3.Code, err3)
	}
	return c.JSON(http.StatusOK, carts)
} 
func (controller cartController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	fmt.Println(id)
	cart, problem := service.Cartservice.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, cart)	
}

func (controller cartController) Update(c echo.Context) error {
		
	cart :=  &model.Cart{}
	if err := c.Bind(cart); err != nil {
		httperror := httperors.NewBadRequestError("Invalid json body")
		return c.JSON(httperror.Code, httperror)
	}	
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	updatedcart, problem := service.Cartservice.Update(id, cart)
	if problem != nil {
		return c.JSON(problem.Code, problem)
	}
	return c.JSON(http.StatusOK, updatedcart)
}

func (controller cartController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code, httperror)
	}
	success, failure := service.Cartservice.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code, failure)
	}
	return c.JSON(success.Code, success)
		
}
