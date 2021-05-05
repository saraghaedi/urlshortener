package router

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/handler"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/log"
	"github.com/sirupsen/logrus"
)

// New creates all the routes.
func New(cfg config.Config, db *gorm.DB) *echo.Echo {
	e := echo.New()

	debug := logrus.IsLevelEnabled(logrus.DebugLevel)

	e.HideBanner = true

	if !debug {
		e.HidePort = true
	}

	e.Debug = debug

	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout

	recoverConfig := middleware.DefaultRecoverConfig
	recoverConfig.DisablePrintStack = !debug
	e.Use(middleware.RecoverWithConfig(recoverConfig))

	e.Use(middleware.CORS())
	e.Use(log.LoggerMiddleware(cfg.Logger.AccessLogger))

	urlRepo := model.SQLURLRepo{
		MasterDB: db,
		SlaveDB:  db,
	}

	urlHandler := handler.NewURLHandler(urlRepo)

	e.POST("/url", urlHandler.NewURL)
	e.GET("/:shortUrl", urlHandler.CallURL)

	return e
}
