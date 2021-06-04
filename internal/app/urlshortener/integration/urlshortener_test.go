package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/response"

	"github.com/jinzhu/gorm"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/cmd"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/request"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/stretchr/testify/suite"
)

const (
	serverHTTPAddress           = "127.0.0.1:7677"
	waitingTimeForDoingSomeJobs = 10 * time.Second
)

type URLShortenerSuite struct {
	suite.Suite
	dbClient *gorm.DB
}

func (suite *URLShortenerSuite) SetupSuite() {
	suite.NoError(os.Setenv("URLSHORTENER_SERVER_ADDRESS", serverHTTPAddress))

	cfg := config.Init()

	dbClient, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	suite.NoError(err)

	suite.dbClient = dbClient

	suite.startServer()
}

func (suite *URLShortenerSuite) TearDownSuite() {
	suite.NoError(os.Unsetenv("URLSHORTENER_SERVER_ADDRESS"))
}

func (suite *URLShortenerSuite) SetupTest() {
	suite.NoError(suite.dbClient.Exec(`truncate table urls`).Error)
}

func (suite *URLShortenerSuite) TearDownTest() {
	suite.NoError(suite.dbClient.Exec(`truncate table urls`).Error)
}

func (suite *URLShortenerSuite) TestScenario() {
	urlRequest := request.NewURL{URL: "https://google.com"}
	urlRequestData, err := json.Marshal(&urlRequest)
	suite.NoError(err)

	req1, err := http.NewRequest("POST", makeURL("/url"), bytes.NewBuffer(urlRequestData))
	suite.NoError(err)
	req1.Header.Set("Content-Type", "application/json")

	client := http.Client{}

	resp1, err := client.Do(req1)
	suite.NoError(err)

	res1Bytes, err := ioutil.ReadAll(resp1.Body)
	suite.NoError(err)

	suite.NoError(resp1.Body.Close())
	suite.Equal(http.StatusOK, resp1.StatusCode)

	var res1Struct response.NewURL

	err = json.Unmarshal(res1Bytes, &res1Struct)
	suite.NoError(err)

	req2, err := http.NewRequest("GET", makeURL(fmt.Sprintf("/%s", res1Struct.ShortURL)), nil)
	suite.NoError(err)

	resp2, err := client.Do(req2)
	suite.NoError(err)

	suite.Equal(http.StatusOK, resp2.StatusCode)
	suite.NoError(resp2.Body.Close())
}

func (suite *URLShortenerSuite) startServer() {
	root := cmd.NewRootCommand()
	if root == nil {
		suite.Fail("root command is nil!")
		return
	}

	root.SetArgs([]string{"server"})

	go func() {
		if err := root.Execute(); err != nil {
			suite.Fail("failed to execute root command: %s", err.Error())
		}
	}()

	// Wait for starting the server
	time.Sleep(waitingTimeForDoingSomeJobs)
}

func makeURL(path string) string {
	return fmt.Sprintf("http://%s%s", serverHTTPAddress, path)
}

func TestURLShortenerSuite(t *testing.T) {
	suite.Run(t, new(URLShortenerSuite))
}
