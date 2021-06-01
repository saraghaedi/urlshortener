package model_test

import (
	"testing"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/stretchr/testify/suite"
)

type URLRepoSuite struct {
	suite.Suite
	repo model.SQLURLRepo
}

func (suite *URLRepoSuite) SetupSuite() {
	cfg := config.Init()
	dbCfg := cfg.Database

	masterDb, err := database.New(dbCfg.Driver, dbCfg.MasterConnStr)
	suite.NoError(err)
	suite.NotNil(masterDb)

	slaveDb, err := database.New(dbCfg.Driver, dbCfg.SlaveConnStr)
	suite.NoError(err)
	suite.NotNil(slaveDb)

	suite.repo = model.SQLURLRepo{
		MasterDB: masterDb,
		SlaveDB:  slaveDb,
	}
}

func (suite *URLRepoSuite) TearDownSuite() {
	suite.NoError(suite.repo.MasterDB.Close())
	suite.NoError(suite.repo.SlaveDB.Close())
}

func (suite *URLRepoSuite) SetupTest() {
	suite.NoError(suite.repo.MasterDB.Exec(`truncate table urls`).Error)
}

func (suite *URLRepoSuite) TearDownTest() {
	suite.NoError(suite.repo.MasterDB.Exec(`truncate table urls`).Error)
}

func (suite *URLRepoSuite) TestCreateAndFind() {
	url := model.URL{
		URL: "github.com/saraghaedi",
	}

	err := suite.repo.Create(&url)
	suite.NoError(err)

	id := url.ID

	findByIDUrl, err := suite.repo.FindByID(uint64(id))
	suite.NoError(err)
	suite.NotNil(findByIDUrl)
	suite.Equal(url.URL, findByIDUrl.URL)
}

func (suite *URLRepoSuite) TestNotFound() {
	findByIDUrl, err := suite.repo.FindByID(10)
	suite.Error(err, model.ErrRecordNotFount)
	suite.Nil(findByIDUrl)
}

func TestURLRepoSuite(t *testing.T) {
	suite.Run(t, new(URLRepoSuite))
}
