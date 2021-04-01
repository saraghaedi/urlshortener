package router

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/saraghaedi/urlshortener/handler"
)

// New creates all the routes.
func New(db *gorm.DB) *echo.Echo {
	e := echo.New()

	urlHandler := handler.NewURLHandler(db)

	e.POST("/url", urlHandler.NewURL)
	e.GET("/:shortUrl", urlHandler.CallURL)

	return e
}
