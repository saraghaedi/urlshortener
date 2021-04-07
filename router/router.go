package router

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saraghaedi/urlshortener/handler"
	"github.com/saraghaedi/urlshortener/log"
	"github.com/sirupsen/logrus"
	"time"
)

// New creates all the routes.
func New(db *gorm.DB) *echo.Echo {
	e := echo.New()

	debug := logrus.IsLevelEnabled(logrus.DebugLevel)

	e.HideBanner = true

	if !debug {
		e.HidePort = true
	}

	e.Debug = debug

	e.Server.ReadTimeout = 20 * time.Second
	e.Server.WriteTimeout = 20 * time.Second

	recoverConfig := middleware.DefaultRecoverConfig
	recoverConfig.DisablePrintStack = !debug
	e.Use(middleware.RecoverWithConfig(recoverConfig))

	e.Use(middleware.CORS())
	e.Use(log.LoggerMiddleware(log.AccessLogger{
		Enabled: true,
		Path:    "./access.log",
	}))

	urlHandler := handler.NewURLHandler(db)

	e.POST("/url", urlHandler.NewURL)
	e.GET("/:shortUrl", urlHandler.CallURL)

	return e
}
