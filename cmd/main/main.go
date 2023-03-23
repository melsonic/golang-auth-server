package main

import (
	handleRoute "example/auth/internal/pkg/controller"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	authorize "example/auth/internal/pkg/authorize"
)

var PORT int = 3000

func main() {

	loaderr := godotenv.Load()

	if loaderr != nil {
		log.Fatal("Error loading .env file")
	}

	router := gin.Default()

	// [x] LOGIN ROUTE
	router.POST("/login", handleRoute.Login)

	// README MAKE /authorize ROUTE TO AUTHORIZE USER BEFORE DOING ANYTHING
	authorizeRoute := router.Group("/authorize")
	authorizeRoute.Use(authorize.ExtractJwtMiddleware())

	// [x] LOGOUT ROUTE
	authorizeRoute.POST("/logout", handleRoute.Logout)

	// [x] VIEW ALL USERS IN THE USER ORGANIZATION
	authorizeRoute.POST("/viewusers", handleRoute.ViewUsers)

	// README MAKE /checkadmin ROUTE TO CHECK IS USER IS ADMIN USER BEFORE MAKING REQUEST
	adminRoute := authorizeRoute.Group("/checkadmin")
	adminRoute.Use(handleRoute.CheckIfAdminUserMiddleware())

	// [x] ADMIN USER ADD NEW USER
	adminRoute.POST("/adduser", handleRoute.AddUser)

	// [x] ADMIN USER DELETE EXISTING USER
	adminRoute.POST("/deleteuser", handleRoute.DeleteUser)

	// app listening at port : 3000
	router.Run(fmt.Sprintf(":%d", PORT))
}
