package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const databaseName = "golinks.db"

// Link is the struct that holds the structure of how information has to be stored
type Link struct {
	gorm.Model
	Short string `form:"short"`
	Url   string `form:"url"`
}

func handleGet(context *gin.Context) {
	short := context.Param("short")
	var link Link

	db, err := gorm.Open("sqlite3", databaseName)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err = db.First(&link, "short = ?", short).Error; gorm.IsRecordNotFoundError(err) {
		context.String(404, "URL not set")
		return
	}
	log.Printf("New Request: %s -> %s", short, link.Url)
	context.Redirect(http.StatusMovedPermanently, link.Url)
}

func newHandler(context *gin.Context) {
	context.HTML(http.StatusOK, "create.tmpl", map[string]string{
		"action": "/new",
	})
}

func newPostHandler(context *gin.Context) {
	var form Link
	db, err := gorm.Open("sqlite3", databaseName)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	if err := context.ShouldBind(&form); err == nil {
		db.Create(&form)
		log.Printf("New Golink created: %s -> %s", form.Short, form.Url)
		context.String(http.StatusOK, "Created")
		return
	}
	log.Fatalf("The link couldn't be created: %v\n", err)
	context.String(http.StatusBadRequest, "Bad Request")
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	// Load up all the templates.
	router.LoadHTMLGlob("templates/*")

	router.GET("/:short", handleGet)
	router.GET("/", newHandler)
	router.POST("/", newPostHandler)

	return router
}

func main() {
	// Setup the database first and migrate the model.
	db, err := gorm.Open("sqlite3", databaseName)
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Link{})
	defer db.Close()
	router := setupRouter()
	router.Run()
}
