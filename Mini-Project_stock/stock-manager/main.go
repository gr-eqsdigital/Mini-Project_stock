package main

// $ compiledaemon --command="./stock-manager"

import (
	"example.com/stock-manager/controllers"
	"example.com/stock-manager/initializers"
	"example.com/stock-manager/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()
	r.POST("/signup", controllers.Signup) // TEST-OK
	r.POST("/login", controllers.Login)   // TEST-OK
	r.POST("/logout", controllers.Logout) // TEST-OK

	var pd = r.Group("product")
	// Set Middlewares
	pd.Use(middleware.RequireAuth)
	pd.Use(middleware.ValidateToken)

	// Register endpoints
	pd.POST("/register", controllers.RegisterProduct)             // TEST-OK
	pd.POST("/update", controllers.UpdateProduct)                 // TEST-FAIL
	pd.POST("/createcategory", controllers.CreateProductCategory) // TEST-OK
	pd.GET("/get/:id", controllers.GetProduct)                    // TEST-OK
	pd.GET("/list", controllers.ListProducts)                     // TEST-OK

	var ug = r.Group("user")

	// Set Middlewares
	ug.Use(middleware.RequireAuth)
	ug.Use(middleware.ValidateToken)

	// Register endpoints
	ug.POST("/update", controllers.UpdateUser) // TEST-OK
	ug.GET("/verify", controllers.VerifyUser)  // TEST-OK
	ug.GET("/get/:email", controllers.GetUser) // TEST-OK
	ug.GET("/list", controllers.ListUsers)     // TEST-OK
	ug.POST("/delete", controllers.DeleteUser) // TEST-OK

	r.Run()
}
