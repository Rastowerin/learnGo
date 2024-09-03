package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Item struct {
	Id          int     `json:"id" gorm:"primary_key"`
	OwnerId     int     `json:"owner_id" binding:"required" gorm:"unique"`
	Label       string  `json:"label" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
}

type ItemControllerAbs interface {
	createItem(c *gin.Context)
	getItems(c *gin.Context)
	getItem(c *gin.Context)
	updateItem(c *gin.Context)
	deleteItem(c *gin.Context)
}

type ItemController struct {
	service ItemServiceAbs
}

func (controller *ItemController) createItem(c *gin.Context) {
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	createdItem, err := controller.service.createItem(item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, createdItem)
}

func (controller *ItemController) getItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	item, err := controller.service.getItem(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (controller *ItemController) getItems(c *gin.Context) {
	items, err := controller.service.getItems()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (controller *ItemController) updateItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var item Item
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	updatedItem, err := controller.service.updateItem(id, item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	c.JSON(http.StatusCreated, updatedItem)
}

func (controller *ItemController) deleteItem(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := controller.service.deleteItem(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	c.JSON(http.StatusNoContent, Item{})
}

func main() {
	ConnectDatabase()

	r := gin.Default()

	r.Use(AuthMiddleware())

	controller := ItemController{
		service: &ItemService{},
	}

	r.POST("/items", controller.createItem)
	r.GET("/items/:id", controller.getItem)
	r.GET("/items", controller.getItems)
	r.PATCH("/items/:id", controller.updateItem)
	r.DELETE("/items/:id", controller.deleteItem)

	er := r.Run("127.0.0.1:8080")
	if er != nil {
		panic(er)
	}
}
