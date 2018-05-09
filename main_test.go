package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

// Test that database is setup properly.
func TestSetupDatabase(t *testing.T) {
	databaseName := "test.db"
	// First, we setup the database.
	setupDatabase(databaseName)
	// Check the sqlite database file was created.
	if _, err := os.Stat(databaseName); os.IsNotExist(err) {
		t.Errorf("Database wasn't created, %s sqlite database doesn't exist.", databaseName)
	}

	// Test that we can open the database that was created.
	db, err := gorm.Open("sqlite3", databaseName)
	if err != nil {
		t.Errorf("Couldn't read the database %s", databaseName)
		return
	}

	// Test that models were migrated.
	db.Create(&Link{Short: "test", Url: "http://example.com"})

	// Test that we can read the database.
	var link Link
	db.First(&link, 1)
	if link.Url != "http://example.com" {
		t.Errorf("Wrong value read from database, expecting http://exmaple.com, got %v", link.Url)
	}

	// Cleanup the database.
	os.Remove(databaseName)

}

// Test Index page handler.
func TestIndexPage(t *testing.T) {
	databaseName := "test.db"
	setupDatabase(databaseName)

	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	os.Remove(databaseName)
}

// Create and test redirection.
func TestRedict(t *testing.T) {
	databaseName := "test.db"
	setupDatabase(databaseName)

	router := setupRouter()
	w := httptest.NewRecorder()

	data := url.Values{}
	data.Set("short", "link")
	data.Add("url", "http://example.com")

	// First, let's create a new URL
	req, _ := http.NewRequest("POST", "/", strings.NewReader(data.Encode()))

	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Created", w.Body.String())

	// Setup a new recorder for redirect.
	w2 := httptest.NewRecorder()
	// Now, let's make sure that redirection also works.
	req, _ = http.NewRequest("GET", "/link", nil)
	router.ServeHTTP(w2, req)

	assert.Equal(t, 301, w2.Code)
	assert.Equal(t, "http://example.com", w2.HeaderMap)

}
