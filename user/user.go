package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct{
	gorm.Model

	Name string `json:"name" blinding:"required"`
}

type UserHandler struct{
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler{
	return &UserHandler{db: db}
}

func (uh *UserHandler)NewUser(c *gin.Context){
	var u User
		err := c.ShouldBindJSON(&u)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),})
			return
		}
		r := uh.db.Create(&u)
		if err = r.Error; err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),})
		}

		c.JSON(http.StatusOK, gin.H{
			"id": u.Model.ID,
			"message": "hello ," + u.Name,})
}

func (uh *UserHandler)GetUser(c *gin.Context) {
	var users []User
	
	r := uh.db.Find(&users)
	if err := r.Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),})
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,})
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can not parse",
	})
	return
}
	var u User
	r := uh.db.Model(User{}).Where("id = ?", idInt).Updates(u)
	if err = r.Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "updated",})
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can not parse",
	})
	return
}
	r := uh.db.Delete(&User{}, idInt)
	if err = r.Error; err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "deleted",})
}
