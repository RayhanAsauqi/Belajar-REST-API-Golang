package main

import (
	"todolist/config"
	"todolist/controllers"

	"github.com/labstack/echo/v4"
)

func main() {
	config.Connect()
	e := echo.New()
	e.GET("/users", controllers.GetUsers)
	e.GET("/users/:id", controllers.GetUser)
	e.POST("/users", controllers.CreateUser)
	e.PATCH("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	e.Logger.Fatal(e.Start(":1234"))
}
