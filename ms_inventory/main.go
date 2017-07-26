package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
	"time"
)

//------------------------------------------ Database Conection
func Database() *gorm.DB {
	//open a db connection
	db, err := gorm.Open("mysql", "root:mysql@/tododb?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {

	//Migrate the schema
	db := Database()
	db.AutoMigrate(&Todo{})
	db.AutoMigrate(&Inventory{})
	db.AutoMigrate(&Product{})

	router := gin.Default()

	inventory := router.Group("/OthalaInventory/ms_i/inventory")
	{
		inventory.POST("/", CreateInventory)
		inventory.GET("/", FetchAllInventory)
		inventory.GET("/:id", FetchSingleInventory)
		inventory.PUT("/:id", UpdateInventory)
		inventory.DELETE("/:id", DeleteInventory)
	}
	products := router.Group("/OthalaInventory/ms_i/product")
	{
		products.POST("/:inventory/", CreateProduct)
		products.GET("/:inventory/", FetchAllProduct)
		//products.GET("/:inventory/:id", FetchSingleProduct)
		//products.PUT("/:inventory/:id", UpdateProduct)
		//products.DELETE("/:inventory/ :id", DeleteProduct)
	}

	v1 := router.Group("/api/v1/todos")
	{
		v1.POST("/", CreateTodo)
		v1.GET("/", FetchAllTodo)
		v1.GET("/:id", FetchSingleTodo)
		v1.PUT("/:id", UpdateTodo)
		v1.DELETE("/:id", DeleteTodo)
	}
	router.Run()

}

//------------------------------------------------- Model -----------------------------------------
type Todo struct {
	gorm.Model
	Title     string    `json:"title"`
	Completed int       `json:"completed"`
	Products  []Product `json:"products"`
}

type Inventory struct {
	gorm.Model
	Products          []Product `json:"products"`
	Name              string    `json:"name"`
	TotalProductValue float64   `json:"products_value"`
	TotalsellsValue   float64   `json:"sells_value"`
	Earnings          float64   `json:"earnings"`
}

type Product struct {
	gorm.Model
	Name        string    `json:"name"`
	Reference   string    `json:"reference"`
	Suplieer    int       `json:"suplieer"`
	EntryDate   time.Time `json:"entry_date"`
	EntryFee    float64   `json:"entry_fee"`
	OutFee      float64   `json:"out_fee"`
	InventoryID uint      `json:"inventory"`
	TodoID      uint      `json:"todo"`
}

type TransformedTodo struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

//-------------------------------------------------- get Products from Inventary------------------------------
func getInventoryProducts(inventoryId uint) []Product {
	var inventory Inventory
	db := Database()
	db.First(&inventory, inventoryId)

	var products []Product
	db.Where("inventory_id = ?", inventoryId).Find(&products)
	fmt.Printf("Products: %v \n", products)
	return products
}

//-------------------------------------------------- API Methods - TODO -----------------------------------
func CreateTodo(c *gin.Context) {
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	todo := Todo{Title: c.PostForm("title"), Completed: completed, Products: []Product{
		Product{
			Name:      "Alice",
			Reference: "A"},
	},
	}
	db := Database()
	db.Save(&todo)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Todo item created successfully!", "resourceId": todo.ID})
}

func FetchAllTodo(c *gin.Context) {
	var todos []Todo
	var _todos []TransformedTodo

	db := Database()
	db.Find(&todos)

	if len(todos) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	//transforms the todos for building a good response
	for _, item := range todos {
		completed := false
		if item.Completed == 1 {
			completed = true
		} else {
			completed = false
		}
		_todos = append(_todos, TransformedTodo{ID: item.ID, Title: item.Title, Completed: completed})
	}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todos})
}

func FetchSingleTodo(c *gin.Context) {
	var todo Todo
	todoId := c.Param("id")

	db := Database()
	db.First(&todo, todoId)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	completed := false
	if todo.Completed == 1 {
		completed = true
	} else {
		completed = false
	}

	_todo := TransformedTodo{ID: todo.ID, Title: todo.Title, Completed: completed}
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": _todo})
}

func UpdateTodo(c *gin.Context) {
	var todo Todo
	todoId := c.Param("id")
	db := Database()
	db.First(&todo, todoId)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Model(&todo).Update("title", c.PostForm("title"))
	completed, _ := strconv.Atoi(c.PostForm("completed"))
	db.Model(&todo).Update("completed", completed)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo updated successfully!"})
}

func DeleteTodo(c *gin.Context) {
	var todo Todo
	todoId := c.Param("id")
	db := Database()
	db.First(&todo, todoId)

	if todo.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No todo found!"})
		return
	}

	db.Delete(&todo)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo deleted successfully!"})
}

//-------------------------------------------------- API Methods - Inventory -----------------------------------
func CreateInventory(c *gin.Context) {
	products_value, _ := strconv.ParseFloat(c.PostForm("products_value"), 64)
	TotalsellsValue, _ := strconv.ParseFloat(c.PostForm("sells_value"), 64)
	Earnings, _ := strconv.ParseFloat(c.PostForm("earnings"), 64)
	var products []Product

	inventory := Inventory{
		Products:          products,
		Name:              c.PostForm("name"),
		TotalProductValue: products_value,
		TotalsellsValue:   TotalsellsValue,
		Earnings:          Earnings,
	}
	db := Database()
	db.Save(&inventory)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Inventory item created successfully!", "resourceId": inventory.ID, "resourceName": inventory.Name})
}

func FetchAllInventory(c *gin.Context) {
	var inventories []Inventory

	db := Database()
	db.Find(&inventories)

	if len(inventories) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Inventory found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": inventories})
}

func FetchSingleInventory(c *gin.Context) {
	var inventory Inventory
	inventoryId := c.Param("id")

	db := Database()
	db.First(&inventory, inventoryId)

	if inventory.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Inventory found!"})
		return
	}

	inventory.Products = getInventoryProducts(inventory.ID)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": inventory})
}

func UpdateInventory(c *gin.Context) {
	var inventory Inventory
	inventoryId := c.Param("id")
	db := Database()
	db.First(&inventory, inventoryId)

	if inventory.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No inventory found!"})
		return
	}

	db.Model(&inventory).Update("name", c.PostForm("name"))
	products_value, _ := strconv.ParseFloat(c.PostForm("products_value"), 64)
	sells_value, _ := strconv.ParseFloat(c.PostForm("sells_value"), 64)
	earnings, _ := strconv.ParseFloat(c.PostForm("earnings"), 64)
	db.Model(&inventory).Update("products_value", products_value)
	db.Model(&inventory).Update("sells_value", sells_value)
	db.Model(&inventory).Update("earnings", earnings)

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Inventory updated successfully!"})
}

func DeleteInventory(c *gin.Context) {
	var inventory Inventory
	inventoryId := c.Param("id")
	db := Database()
	db.First(&inventory, inventoryId)

	if inventory.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No inventory found!"})
		return
	}

	db.Delete(&inventory)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Todo inventory successfully!"})
}

//-------------------------------------------------- API Methods - Product -----------------------------------
func CreateProduct(c *gin.Context) {
	var inventory Inventory
	inventoryId := c.Param("inventory")
	db := Database()
	db.First(&inventory, inventoryId)

	suplieer, _ := strconv.Atoi(c.PostForm("suplieer"))
	entry_date, _ := time.Parse(time.RFC3339, c.PostForm("entry_date"))
	entry_fee, _ := strconv.ParseFloat(c.PostForm("entry_fee"), 64)
	out_fee, _ := strconv.ParseFloat(c.PostForm("out_fee"), 64)

	product := Product{
		Name:        c.PostForm("name"),
		Reference:   c.PostForm("reference"),
		Suplieer:    suplieer,
		EntryDate:   entry_date,
		EntryFee:    entry_fee,
		OutFee:      out_fee,
		InventoryID: inventory.ID,
	}
	db.Save(&product)
	c.JSON(http.StatusCreated, gin.H{"status": http.StatusCreated, "message": "Product item created successfully!", "resourceId": product.ID, "Inventory Products": inventory.Products})
}

func FetchAllProduct(c *gin.Context) {
	var products []Product

	db := Database()
	db.Find(&products)

	if len(products) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No Products found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": products})
}
