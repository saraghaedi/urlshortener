package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/labstack/echo/v4"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/request"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/response"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/utils"
)

// URLHandler handles all incoming requests related to URLs.
type URLHandler struct {
	URLRepo model.URLRepo
}

// NewURLHandler returns a new URLHandler.
func NewURLHandler(urlRepo model.URLRepo) *URLHandler {
	return &URLHandler{
		URLRepo: urlRepo,
	}
}

// NewURL creates a new short URL.
func (u URLHandler) NewURL(c echo.Context) error {
	var req request.NewURL

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("failed to decode request: %s", err.Error())
		return echo.ErrBadRequest
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("Create URL validation : %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	url := &model.URL{URL: req.URL}

	if err := u.URLRepo.Create(url); err != nil {
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

	url, err := u.URLRepo.FindByID(id)
	if err != nil {
		if err == model.ErrRecordNotFound {
			return echo.ErrNotFound
		}

		logrus.Errorf("failed to read url from db: %s", err.Error())

		return echo.ErrInternalServerError
	}

	return c.Redirect(http.StatusPermanentRedirect, url.URL)
}
