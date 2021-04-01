package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/saraghaedi/urlshortener/request"
	"github.com/saraghaedi/urlshortener/response"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/functions"
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

// NewURL will add a new url in database.
func (u URLHandler) NewURL(c echo.Context) error {
	var req request.NewURL

	_ = c.Bind(&req)

	url := &model.URL{URL: req.URL}

	u.db.Create(url)

	id := url.ID

	shortURL := functions.URLEncoder(int64(id))

	resp := response.NewURL{ShortURL: shortURL}

	return c.JSON(http.StatusOK, resp)
}
