package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

func prepareItem() Item {
	return Item{
		OwnerId:     1,
		Label:       "mashina",
		Description: "bystraya",
		Price:       5.0,
	}
}

func setupRouter() *gin.Engine {
	testItemService := &ItemService{}
	testItemController := ItemController{
		service: testItemService,
	}

	testItemService.createItem(
		Item{
			OwnerId:     2,
			Label:       "mashina",
			Description: "bystraya",
			Price:       5.0,
		})

	r := gin.Default()
	r.POST("/items", testItemController.createItem)
	r.GET("/items/:id", testItemController.getItem)
	r.GET("/items", testItemController.getItem)
	r.PATCH("/items/:id", testItemController.updateItem)
	r.DELETE("/items/:id", testItemController.deleteItem)
	return r
}

func TestMain(m *testing.M) {

	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	err = DB.AutoMigrate(&Item{})
	if err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	code := m.Run()

	os.Exit(code)
}

func TestCreateItem(t *testing.T) {
	r := setupRouter()

	jsonValue, _ := json.Marshal(prepareItem())
	req, _ := http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestGetItem(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/items/1", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	item := prepareItem()
	item.Id = 1
	item.OwnerId = 2
	itemStr, _ := json.Marshal(item)
	assert.JSONEqf(t, string(itemStr), w.Body.String(), "error")
}

func TestGetItems(t *testing.T) {
	r := setupRouter()

	req, _ := http.NewRequest("GET", "/items", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateItem(t *testing.T) {
	r := setupRouter()

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/items/"+strconv.Itoa(prepareItem().Id), nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	item := prepareItem()
	itemStr, _ := json.Marshal(item)
	assert.JSONEqf(t, string(itemStr), w.Body.String(), "error")
}

func TestDeleteItem(t *testing.T) {
	r := setupRouter()

	jsonValue, _ := json.Marshal(prepareItem())
	_, _ = http.NewRequest("POST", "/items", bytes.NewBuffer(jsonValue))

	req, _ := http.NewRequest("DELETE", "/items/"+strconv.Itoa(prepareItem().Id), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
