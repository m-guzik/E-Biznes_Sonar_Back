package main

import (
	"net/http"
    "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

type Payment struct {
    gorm.Model
    Name string
    Surname string
    Phone string
    Email string
    Address string
    CreditCardNumber string
    CreditCardExpirationDate string
    CreditCardCVV string
    Sum uint
}

type Product struct {
    gorm.Model
    Name string
    Price uint
}

type CartItem struct {
    gorm.Model
    ProductId uint
    Name string
    Price uint
    Amount uint
}

type PaymentController struct {
    db *gorm.DB
}

type ProductController struct {
    db *gorm.DB
}

type CartController struct {
    db *gorm.DB
}


func (paymentController *PaymentController) makePayment(c echo.Context) error {
    var payment Payment
    c.Bind(&payment)
    paymentController.db.Create(&payment)
    return c.JSON(201, payment)
}

func (paymentController *PaymentController) showAllPayments(c echo.Context) error {
    var payments []Payment
    paymentController.db.Find(&payments)
    return c.JSON(200, payments)
}

func (paymentController *PaymentController) deletePayment(c echo.Context) error {
    id := c.Param("id")
    var payment Payment
    paymentController.db.Delete(&payment, id)
    return c.JSON(200, "")
}


func (productController *ProductController) addProduct(c echo.Context) error {
    var product Product
    c.Bind(&product)
    productController.db.Create(&product)
    return c.JSON(201, product)
}

func (productController *ProductController) showAllProducts(c echo.Context) error {
    var products []Product
    productController.db.Find(&products)
    return c.JSON(200, products)
}

func (productController *ProductController) showProduct(c echo.Context) error {
    id := c.Param("id")
    var product Product
    productController.db.Take(&product, id)
    return c.JSON(200, product)
}

func (productController *ProductController) deleteProduct(c echo.Context) error {
    id := c.Param("id")
    var product Product
    productController.db.Delete(&product, id)
    return c.JSON(200, "")
}

func (productController *ProductController) updatePrice(c echo.Context) error {
    id := c.Param("id")
    newPrice := c.Param("price")
    var product Product
    productController.db.Take(&product, id)
    productController.db.Model(&product).Update("Price", newPrice)
    return c.JSON(200, product)
}


func (cartController *CartController) addToCart(c echo.Context) error {
    var cartItem CartItem

	var request struct {
		ProductID uint `json:"ID"`
		ProductName string `json:"Name"`
		ProductPrice uint `json:"Price"`
	}
	c.Bind(&request);

    cartController.db.Find(&cartItem, CartItem{ProductId: request.ProductID})
    currentAmount := cartItem.Amount

    if currentAmount > 0 {
    	cartController.db.Model(&CartItem{}).Where("product_id = ?", request.ProductID).Update("Amount", currentAmount + 1)
    } else {
    	cartItem := CartItem{
            ProductId: request.ProductID,
            Name: request.ProductName,
            Price: request.ProductPrice,
            Amount: 1,
    	}
    	cartController.db.Create(&cartItem)
    }
    return c.JSON(201, cartItem)
}

func (cartController *CartController) totalCartValue(c echo.Context) error {
    var sum int
    cartController.db.Table("cart_items").Select("sum(price*amount)").Row().Scan(&sum)
    return c.JSON(200, sum)
}

func (cartController *CartController) showCart(c echo.Context) error {
    var cartItems []CartItem
    cartController.db.Find(&cartItems)
    return c.JSON(200, cartItems)
}

func (cartController *CartController) deleteFromCart(c echo.Context) error {
    id := c.Param("id")
    var cartItem CartItem
    cartController.db.Unscoped().Delete(&cartItem, id)
    return c.JSON(200, "")
}


func main() {
    db, err := gorm.Open(sqlite.Open("products.db"), &gorm.Config{})
    if err != nil {
        panic("Failed to connect database")
    }

    db.AutoMigrate(&Product{})
    db.AutoMigrate(&Payment{})
    db.AutoMigrate(&CartItem{})

    productController := &ProductController{db : db}
    paymentController := &PaymentController{db : db}
    cartController := &CartController{db : db}

	e := echo.New()

    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{"*"},
        AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
        AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
    }))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello and welcome to my Bike Shop!")
	})

	e.POST("addProduct", productController.addProduct)
	e.GET("products", productController.showAllProducts)
	e.GET("product/:id", productController.showProduct)
	e.DELETE("delete/:id", productController.deleteProduct)
	e.PUT("changePrice/:id/:price", productController.updatePrice)

    e.GET("cart", cartController.showCart)
    e.GET("totalCartValue", cartController.totalCartValue)
    e.POST("addToCart", cartController.addToCart)
    e.DELETE("deleteFromCart/:id", cartController.deleteFromCart)

	e.POST("makePayment", paymentController.makePayment)
	e.GET("payments", paymentController.showAllPayments)
    e.DELETE("deletePayment/:id", paymentController.deletePayment)

	e.Logger.Fatal(e.Start(":8070"))
}

