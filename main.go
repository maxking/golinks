package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var databaseName = "golinks.db"

// Link represents a data structure that will be used to store the mapping from large to short URLs.
type Link struct {
	gorm.Model
	Short string `form:"short"`
	Url   string `form:"url"`
}

func setupDatabase(databaseName string) {
	// Setup the database first and migrate the model.
	db, err := gorm.Open("sqlite3", databaseName)
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Link{})
	db.Close()
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

	if err := context.Bind(&form); err == nil {
		if form.Short == "" || form.Url == "" {
			// Empty short link.
			context.String(http.StatusBadRequest,
				fmt.Sprintf("Bad parameters short=%s url=%s", form.Short, form.Url))
			return
		}

		// If the short is already set in the database, return an error.
		var link Link
		err = db.Where("short = ?", form.Short).First(&link).Error
		if err == nil {
			context.String(http.StatusConflict, fmt.Sprintf("URL already exists and points to %s", link.Url))
			return
		}

		db.Create(&form)
		log.Printf("New Golink created: %s -> %s", form.Short, form.Url)
		context.String(http.StatusOK, "Created")
		return
	}

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
	setupDatabase(databaseName)

	router := setupRouter()
	router.Run()
}
