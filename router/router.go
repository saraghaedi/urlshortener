package router

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/handler"
)

// New handles all the routes.
func New(db *gorm.DB) *mux.Router {
	r := mux.NewRouter().StrictSlash(true)

	urlHandler := handler.NewURLHandler(db)

	r.HandleFunc("/urls", urlHandler.AllURLs).Methods("GET")
	//r.HandleFunc("/url", urlHandler.URLShortener).Methods("GET")
	r.HandleFunc("/url/{url}", urlHandler.NewURL).Methods("POST")
	r.HandleFunc("/url/{id}", urlHandler.DeleteURL).Methods("DELETE")
	r.HandleFunc("/url/{id}/{url}", urlHandler.UpdateURL).Methods("PUT")

	return r
}
