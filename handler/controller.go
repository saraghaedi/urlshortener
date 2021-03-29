package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/model"
)

// URLHandler will hold everything that controller needs.
type URLHandler struct {
	db *gorm.DB
}

// NewURLHandler returns a new BaseHandler.
func NewURLHandler(db *gorm.DB) *URLHandler {
	return &URLHandler{
		db: db,
	}
}

// AllURLs returns all Urls in database.
func (u URLHandler) AllURLs(w http.ResponseWriter, r *http.Request) {
	var urls []model.URL

	u.db.Find(&urls)

	err := json.NewEncoder(w).Encode(urls)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// NewURL will add a new url in database.
func (u URLHandler) NewURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bigurl := vars["url"]

	u.db.Create(&model.URL{URL: bigurl})

	_, err := fmt.Fprintf(w, "New URL added to the database")

	if err != nil {
		fmt.Println(err.Error())
	}
}

//DeleteURL delete a specific url with id from the database.
func (u URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idparam := vars["id"]

	u.db.Where("id = ?", idparam).Delete(model.URL{})

	_, err := fmt.Fprintf(w, "Successfully Deleted URL")

	if err != nil {
		fmt.Println(err.Error())
	}
}

// UpdateURL update a url by id from the database.
func (u URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idparam := vars["id"]
	newURL := vars["url"]

	u.db.Where("id = ?", idparam).Update("url", &newURL)

	_, err := fmt.Fprintf(w, "Successfully Updated URL")

	if err != nil {
		fmt.Println(err.Error())
	}
}
