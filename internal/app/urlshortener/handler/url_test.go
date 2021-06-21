package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/response"

	"github.com/labstack/echo/v4"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/handler"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/request"
	"github.com/stretchr/testify/suite"
)

type fakeURLRepo struct {
	model.URLRepo
	repoError error
}

type fakeURLCounterRepo struct {
	model.URLCounterRepo
	repoError error
}

func (f *fakeURLRepo) Create(url *model.URL) error {
	url.ID = 1
	return f.repoError
}

func (f *fakeURLRepo) FindByID(id uint64) (*model.URL, error) {
	if f.repoError != nil {
		return nil, f.repoError
	}

	url := &model.URL{
		Model: gorm.Model{
			ID: 1,
		},
		URL: "https://github.com/saraghaedi",
	}

	if id == uint64(url.ID) {
		return url, nil
	}

	return nil, model.ErrRecordNotFound
}

func (f *fakeURLCounterRepo) Incr(id uint64) error {
	return f.repoError
}

type URLHandlerSuite struct {
	suite.Suite
	engine             *echo.Echo
	fakeURLRepo        *fakeURLRepo
	fakeURLCounterRepo *fakeURLCounterRepo
}

func (suite *URLHandlerSuite) SetupSuite() {
	suite.engine = echo.New()

	suite.fakeURLRepo = &fakeURLRepo{}
	suite.fakeURLCounterRepo = &fakeURLCounterRepo{}

	urlHandler := handler.URLHandler{
		URLRepo:        suite.fakeURLRepo,
		URLCounterRepo: suite.fakeURLCounterRepo,
	}

	suite.engine.POST("/url", urlHandler.NewURL)
	suite.engine.GET("/:shortUrl", urlHandler.CallURL)
}

func (suite *URLHandlerSuite) TestNewURL() {
	cases := []struct {
		name      string
		req       request.NewURL
		status    int
		repoError error
	}{
		{
			name:      "successfully create a new shorted URL",
			req:       request.NewURL{URL: "github.com/saraghaedi"},
			status:    http.StatusOK,
			repoError: nil,
		},
		{
			name:      "Failed to create shorted URL: Empty URL",
			req:       request.NewURL{URL: ""},
			status:    http.StatusBadRequest,
			repoError: nil,
		},
		{
			name:      "Failed to create shorted URL: Wrong URL format",
			req:       request.NewURL{URL: "Hello"},
			status:    http.StatusBadRequest,
			repoError: nil,
		},
		{
			name:      "Failed to create shorted URL: Repo error",
			req:       request.NewURL{URL: "github.com/saraghaedi"},
			status:    http.StatusInternalServerError,
			repoError: errors.New("something went wrong"),
		},
	}
	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeURLRepo.repoError = tc.repoError

			data, err := json.Marshal(tc.req)
			suite.NoError(err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/url", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)

			if tc.status == http.StatusOK {
				var resp response.NewURL

				suite.NoError(json.Unmarshal(w.Body.Bytes(), &resp))

				suite.NotEmpty(resp.ShortURL)
			}
		})
	}
}

func (suite *URLHandlerSuite) TestCallURL() {
	cases := []struct {
		name           string
		shortedURL     string
		status         int
		sqlRepoError   error
		redisRepoError error
		resp           response.NewURL
	}{
		{
			name:           "successfully call shorted URL ",
			shortedURL:     "1",
			status:         http.StatusTemporaryRedirect,
			sqlRepoError:   nil,
			redisRepoError: nil,
			resp:           response.NewURL{ShortURL: "1"},
		},
		{
			name:           "failed to call shorted URL",
			shortedURL:     "10",
			status:         http.StatusNotFound,
			sqlRepoError:   nil,
			redisRepoError: nil,
		},
		{
			name:           "sql repo error",
			shortedURL:     "1",
			status:         http.StatusInternalServerError,
			sqlRepoError:   errors.New("something went wrong"),
			redisRepoError: nil,
		},
		{
			name:           "redis repo error",
			shortedURL:     "1",
			status:         http.StatusTemporaryRedirect,
			redisRepoError: errors.New("something went wrong"),
			sqlRepoError:   nil,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeURLRepo.repoError = tc.sqlRepoError
			suite.fakeURLCounterRepo.repoError = tc.redisRepoError

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/%s", tc.shortedURL), nil)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)
		})
	}
}

func TestURLHandlerSuite(t *testing.T) {
	suite.Run(t, new(URLHandlerSuite))
}
