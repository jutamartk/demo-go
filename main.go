package main

import (
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/jutamartk/demo-gin/auth"
	"github.com/jutamartk/demo-gin/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil{
		log.Fatal("Can not Load env!")
	}
	
	dsn := os.Getenv("DATABASE")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		log.Printf("can not connect DB!")
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		//response
		// c.String(http.StatusOK,"pong")

		c.JSON(http.StatusOK, gin.H{
			"message": "pong!",
		})
	})

	r.GET("/greeting",func(c *gin.Context) {
		//queryParam
		name := c.Query("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong!" + name,
		})
	})

	r.GET("/hello/:name",func(c *gin.Context) {
		name := c.Query("name")
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello!" + name,
		})
	})


	db.AutoMigrate(
		&user.User{},
	)
	r.GET("/tokenz", auth.GetToken(os.Getenv("WIFI_SECRET")))
	//Route ต่อจากนี้จะติดmiddleware ที่ต้องมี token
	authRoute := r.Group("", auth.AuthMiddleware(os.Getenv("WIFI_SECRET")))

	userHandler := user.NewUserHandler(db)
	authRoute.POST("/users", userHandler.NewUser)

	authRoute.GET("/users",userHandler.GetUser )
	//PATCH
	authRoute.PATCH("/users/:id", userHandler.UpdateUser)
		
	// DELETE
	authRoute.DELETE("/users/:id",userHandler.DeleteUser)

	r.Run()
}