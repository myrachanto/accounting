package routes

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/myrachanto/accounting/controllers"
	jwt "github.com/dgrijalva/jwt-go"
)

func StoreApi() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file in routes")
	}
	PORT := os.Getenv("PORT")
	key := os.Getenv("EncryptionKey")

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover()) 
	e.Use(middleware.CORS())

	e.Static("/", "public")
	JWTgroup := e.Group("/api/")
	JWTgroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey: []byte(key),
	}))
// e.Set(id, 1)

	// // Get retrieves data from the context.
	// Get(key string) interface{}

	// // Set saves data in the context.
	// Set(key string, val interface{}) 

	// admin := e.Group("admin/")
	// admin.Use(isAdmin)
	/////////////////////////////////////////////////////////////////////////////////////
	////////////////////////needs more info ////////////////////////////////////////////
	///////////////////////////////////////////////////////////////////////////////////
	// var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	// 	SigningMethod: "HS256",
	// 	SigningKey: []byte(key),
	// })
	//JwtG := e.Group("/users")
	// JwtG.Use(middleware.JWT([]byte(key)))
	// Routes
	//e.GET("/is-loggedin", h.private, IsLoggedIn)
	//e.POST("/register", IsLoggedIn,isAdmin,isEmployee,isSupervisor, controllers.UserController.Create)
	e.POST("/register", controllers.UserController.Create)
	e.POST("/login", controllers.UserController.Login)
	JWTgroup.GET("logout/:token", controllers.UserController.Logout)
	JWTgroup.GET("users", controllers.UserController.GetAll)
	JWTgroup.GET("users/:id", controllers.UserController.GetOne)
	JWTgroup.PUT("users/:id", controllers.UserController.Update)
	JWTgroup.DELETE("users/:id", controllers.UserController.Delete)

	///////////dashboard/////////////////////////////	
	JWTgroup.GET("dashboard", controllers.DashboardController.View)
	JWTgroup.GET("email/create", controllers.DashboardController.Email)
	// JWTgroup.POST("email/create", controllers.DashboardController.Send)
	//e.DELETE("loggoutall/:id", controllers.UserController.DeleteALL) logout all accounts
	///////////message/////////////////////////////	
	JWTgroup.POST("messages", controllers.MessageController.Create)
	JWTgroup.GET("messages", controllers.MessageController.GetAll)
	JWTgroup.GET("messages/:id", controllers.MessageController.GetOne)
	JWTgroup.PUT("messages/:id", controllers.MessageController.Update)
	JWTgroup.DELETE("messages/:id", controllers.MessageController.Delete)
	///////////nortifications/////////////////////////////	
	JWTgroup.POST("nortifications", controllers.NortificationController.Create)
	JWTgroup.GET("nortifications", controllers.MessageController.GetAll) 
	JWTgroup.GET("nortifications/:id", controllers.MessageController.GetOne)
	JWTgroup.PUT("nortifications/:id", controllers.MessageController.Update)
	JWTgroup.DELETE("nortifications/:id", controllers.MessageController.Delete)
	///////////category/////////////////////////////	
	JWTgroup.GET("categorys/view", controllers.CategoryController.View)
	JWTgroup.POST("categorys", controllers.CategoryController.Create)
	JWTgroup.GET("categorys", controllers.CategoryController.GetAll)
	JWTgroup.GET("categorys/:id", controllers.CategoryController.GetOne)
	JWTgroup.PUT("categorys/:id", controllers.CategoryController.Update)
	JWTgroup.DELETE("categorys/:id", controllers.CategoryController.Delete)
	///////////majorcategory/////////////////////////////	
	JWTgroup.POST("majorcategory", controllers.MCategoryController.Create)
	JWTgroup.GET("majorcategory", controllers.MCategoryController.GetAll)
	JWTgroup.GET("majorcategory/:id", controllers.MCategoryController.GetOne)
	JWTgroup.PUT("majorcategory/:id", controllers.MCategoryController.Update)
	JWTgroup.DELETE("majorcategory/:id", controllers.MCategoryController.Delete)
	///////////subcategory/////////////////////////////	
	JWTgroup.POST("subcategory", controllers.SubcategoryController.Create)
	JWTgroup.GET("subcategory", controllers.SubcategoryController.GetAll)
	JWTgroup.GET("subcategory/:id", controllers.SubcategoryController.GetOne)
	JWTgroup.PUT("subcategory/:id", controllers.SubcategoryController.Update)
	JWTgroup.DELETE("subcategory/:id", controllers.SubcategoryController.Delete)
	///////////subcategory/////////////////////////////	
	JWTgroup.GET("products/view", controllers.ProductController.View)
	JWTgroup.POST("products", controllers.ProductController.Create)
	e.GET("productsearch", controllers.ProductController.GetProducts)
	JWTgroup.GET("products", controllers.ProductController.GetAll)
	e.GET("products/:id", controllers.ProductController.GetOne)
	JWTgroup.GET("products/:id", controllers.ProductController.GetOne)
	JWTgroup.PUT("products/:id", controllers.ProductController.Update)
	JWTgroup.DELETE("products/:id", controllers.ProductController.Delete)
	///////////cart/////////////////////////////	
	// JWTgroup.GET("cart/create", controllers.CartController.View)
	JWTgroup.POST("cart", controllers.CartController.Create)
	JWTgroup.GET("cart", controllers.CartController.GetAll)
	JWTgroup.GET("cart/:id", controllers.CartController.GetOne)
	JWTgroup.PUT("cart/:id", controllers.CartController.Update)
	JWTgroup.DELETE("cart/:id", controllers.CartController.Delete)
	//////////////////////////////////////////////////////////////////////////
	///////////////////customer module///////////////////////////////////////
	///////////Invoice/////////////////////////////	////////////////////////
	JWTgroup.POST("customer", controllers.CustomerController.Create)
	JWTgroup.GET("customer", controllers.CustomerController.GetAll)
	JWTgroup.GET("customer/:id", controllers.CustomerController.GetOne)
	JWTgroup.PUT("customer/:id", controllers.CustomerController.Update)
	JWTgroup.DELETE("customer/:id", controllers.CustomerController.Delete)
	///////////Invoice/////////////////////////////	
	JWTgroup.GET("invoice/create", controllers.InvoiceController.View)
	JWTgroup.POST("invoice", controllers.InvoiceController.Create)
	JWTgroup.GET("invoice", controllers.InvoiceController.GetAll)
	JWTgroup.GET("invoice/:id", controllers.InvoiceController.GetOne)
	// JWTgroup.PUT("invoice/:id", controllers.InvoiceController.Update)
	// JWTgroup.DELETE("invoice/:id", controllers.InvoiceController.Delete)
	///////////trasanctions/////////////////////////////	
	JWTgroup.POST("trasanctions", controllers.TransactionController.Create)
	JWTgroup.GET("trasanctions", controllers.TransactionController.GetAll)
	JWTgroup.GET("trasanctions/:id", controllers.TransactionController.GetOne)
	JWTgroup.PUT("trasanctions/:id", controllers.TransactionController.Update)
	JWTgroup.DELETE("trasanctions/:id", controllers.TransactionController.Delete)
	//////////////////////////////////////////////////////////////////////////
	///////////////////supplier module///////////////////////////////////////
	///////////Invoice/////////////////////////////	//////////////////////////
	JWTgroup.POST("supplier", controllers.SupplierController.Create)
	JWTgroup.GET("supplier", controllers.SupplierController.GetAll)
	JWTgroup.GET("supplier/:id", controllers.SupplierController.GetOne)
	JWTgroup.PUT("supplier/:id", controllers.SupplierController.Update)
	JWTgroup.DELETE("supplier/:id", controllers.SupplierController.Delete)
	///////////Invoice/////////////////////////////	
	JWTgroup.GET("Viewsinvoice", controllers.SinvoiceController.View)
	JWTgroup.POST("sinvoice", controllers.SinvoiceController.Create)
	JWTgroup.GET("sinvoice", controllers.SinvoiceController.GetAll)
	JWTgroup.GET("sinvoice/:id", controllers.SinvoiceController.GetOne)
	// JWTgroup.PUT("sinvoice/:id", controllers.SinvoiceController.Update)
	// JWTgroup.DELETE("sinvoice/:id", controllers.SinvoiceController.Delete)
	///////////trasanctions/////////////////////////////	
	JWTgroup.POST("strasanctions", controllers.StransactionController.Create)
	JWTgroup.GET("strasanctions", controllers.StransactionController.GetAll)
	JWTgroup.GET("strasanctions/:id", controllers.StransactionController.GetOne)
	JWTgroup.PUT("strasanctions/:id", controllers.StransactionController.Update)
	JWTgroup.DELETE("strasanctions/:id", controllers.StransactionController.Delete)
	//////////////////////////////////////////////////////////////////////////
	///////////////////finance module///////////////////////////////////////
	///////////payments/////////////////////////////////////////////////////
	JWTgroup.POST("payments", controllers.PaymentController.Create)
	JWTgroup.GET("payments", controllers.PaymentController.GetAll)
	JWTgroup.GET("payments/:id", controllers.PaymentController.GetOne)
	// JWTgroup.PUT("payments/:id", controllers.PaymentController.Update)
	// JWTgroup.DELETE("payments/:id", controllers.PaymentController.Delete)
	///////////receipts/////////////////////////////////////////////////////
	JWTgroup.POST("receipts", controllers.ReceiptController.Create)
	JWTgroup.GET("receipts", controllers.ReceiptController.GetAll)
	JWTgroup.GET("receipts/:id", controllers.ReceiptController.GetOne)
	// JWTgroup.PUT("receipts/:id", controllers.ReceiptController.Update)
	// JWTgroup.DELETE("receipts/:id", controllers.ReceiptController.Delete)
	///////////payrecpt/////////////////////////////////////////////////////
	JWTgroup.GET("Viewspayrecpt", controllers.PayrectrasanController.View)
	JWTgroup.POST("payrecpt", controllers.PayrectrasanController.Create)
	JWTgroup.GET("payrecpt", controllers.PayrectrasanController.GetAll)
	JWTgroup.GET("payrecpt/:id", controllers.PayrectrasanController.GetOne)
	///////////Assets/////////////////////////////////////////////////////
	JWTgroup.POST("assets", controllers.AssetController.Create)
	JWTgroup.GET("assets", controllers.AssetController.GetAll)
	JWTgroup.GET("assets/:id", controllers.AssetController.GetOne)
	///////////Assets/////////////////////////////////////////////////////
	JWTgroup.POST("assetstransactions", controllers.AsstransController.Create)
	JWTgroup.GET("assetstransactions", controllers.AsstransController.GetAll)
	JWTgroup.GET("assetstransactions/:id", controllers.AsstransController.GetOne)
	///////////Assets///////////////////////////////////////////////////// 
	JWTgroup.POST("liability", controllers.LiabilityController.Create)
	JWTgroup.GET("liability", controllers.LiabilityController.GetAll)
	JWTgroup.GET("liability/:id", controllers.LiabilityController.GetOne)
	///////////Assets/////////////////////////////////////////////////////
	JWTgroup.POST("liatransanctions", controllers.LiatranController.Create)
	JWTgroup.GET("liatransanctions", controllers.LiatranController.GetAll)
	JWTgroup.GET("liatransanctions/:id", controllers.LiatranController.GetOne)
	///////////Expence/////////////////////////////////////////////////////
	JWTgroup.POST("expence", controllers.ExpenceController.Create)
	JWTgroup.GET("expence", controllers.ExpenceController.GetAll)
	JWTgroup.GET("expence/:id", controllers.ExpenceController.GetOne)
	///////////expencetans/////////////////////////////////////////////////////
	JWTgroup.POST("expencetransanctions", controllers.ExpencetrasanController.Create)
	JWTgroup.GET("expencetransanctions", controllers.ExpencetrasanController.GetAll)
	JWTgroup.GET("expencetransanctions/:id", controllers.ExpencetrasanController.GetOne)
	//////////////////////////////////////////////////////////////////////////
	///////////////////Miscellenous module///////////////////////////////////////
	///////////prices/////////////////////////////////////////////////////
	JWTgroup.POST("prices", controllers.PriceController.Create)
	JWTgroup.GET("prices", controllers.PriceController.GetAll)
	JWTgroup.GET("prices/:id", controllers.PriceController.GetOne)
	JWTgroup.PUT("prices/:id", controllers.PriceController.Update)
	JWTgroup.DELETE("prices/:id", controllers.PriceController.Delete)
	///////////tax/////////////////////////////////////////////////////
	JWTgroup.POST("tax", controllers.TaxController.Create)
	JWTgroup.GET("tax", controllers.TaxController.GetAll)
	JWTgroup.GET("tax/:id", controllers.TaxController.GetOne)
	JWTgroup.PUT("tax/:id", controllers.TaxController.Update)
	JWTgroup.DELETE("tax/:id", controllers.TaxController.Delete)
	///////////discounts//////////////////////////////////////////
	JWTgroup.POST("discounts", controllers.DiscountController.Create)
	JWTgroup.GET("discounts", controllers.DiscountController.GetAll)
	JWTgroup.GET("discounts/:id", controllers.DiscountController.GetOne)
	JWTgroup.PUT("discounts/:id", controllers.DiscountController.Update)
	JWTgroup.DELETE("discounts/:id", controllers.DiscountController.Delete)
	///////////scart/////////////////////////////	
	JWTgroup.GET("Viewscart", controllers.ScartController.View)
	JWTgroup.POST("scart", controllers.ScartController.Create)
	JWTgroup.GET("scart", controllers.ScartController.GetAll)
	JWTgroup.GET("scart/:id", controllers.ScartController.GetOne)
	JWTgroup.PUT("scart/:id", controllers.ScartController.Update)
	JWTgroup.DELETE("scart/:id", controllers.ScartController.Delete)

	// Start server
	e.Logger.Fatal(e.Start(PORT))
}
func isAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isAdmin := claims["Admin"].(bool)
		if isAdmin == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func isSupervisor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isSupervisor := claims["Supervisor"].(bool)
		if isSupervisor == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}
func isEmployee(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("uname").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		isEmployee := claims["Employee"].(bool)
		if isEmployee == false {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}