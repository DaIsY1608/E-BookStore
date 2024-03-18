package main

import (
	"test/dilaf/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.GET("/allbooks", controllers.Allbooks)
	r.Run(":2020")

}
