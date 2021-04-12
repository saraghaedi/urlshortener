package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/saraghaedi/urlshortener/request"
	"github.com/saraghaedi/urlshortener/response"

	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/model"
	"github.com/saraghaedi/urlshortener/utils"
)

// URLHandler handles all incoming requests related to URLs.
type URLHandler struct {
	db *gorm.DB
}

// NewURLHandler returns a new URLHandler.
func NewURLHandler(db *gorm.DB) *URLHandler {
	return &URLHandler{
		db: db,
	}
}

// NewURL creates a new short URL.
func (u URLHandler) NewURL(c echo.Context) error {
	var req request.NewURL

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("failed to decode request: %s", err.Error())
		return echo.ErrBadRequest
	}

	url := &model.URL{URL: req.URL}

	if err := u.db.Create(url).Error; err != nil {
		logrus.Errorf("failed to create new url: %s", err.Error())
		return echo.ErrInternalServerError
	}

	id := url.ID

	shortURL := utils.Base36Encoder(int64(id))

	resp := response.NewURL{ShortURL: shortURL}

	return c.JSON(http.StatusOK, resp)
}

// CallURL redirects to the main URL.
func (u URLHandler) CallURL(c echo.Context) error {
	shortURL := c.Param("shortUrl")

	id := utils.Base36Decoder(shortURL)

	var url model.URL

	if err := u.db.Where("id = ?", id).Take(&url).Error; err != nil {
		logrus.Errorf("failed to read url from db: %s", err.Error())
		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusMovedPermanently, url.URL)
}
