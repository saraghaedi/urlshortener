package model_test

import (
	"testing"

	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type RedisURLRepoSuite struct {
	suite.Suite
	RedisRepo model.RedisURLRepo
}

func (suite *RedisURLRepoSuite) SetupSuite() {
	cfg := config.Init()

	masterDb, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to master database: %s", err.Error())
	}

	slaveDb, err := database.New(cfg.Database.Driver, cfg.Database.SlaveConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to slave database: %s", err.Error())
	}

	redisCfg := cfg.Redis

	rMaster, _ := redis.Create(redisCfg.MasterAddress, redisCfg.Options, true)
	rSlave, _ := redis.Create(redisCfg.SlaveAddress, redisCfg.Options, false)

	suite.RedisRepo = model.RedisURLRepo{
		Base: model.SQLURLRepo{
			MasterDB: masterDb,
			SlaveDB:  slaveDb,
		},
		RedisMasterClient: rMaster,
		RedisSlaveClient:  rSlave,
	}
}

func (suite *RedisURLRepoSuite) SetupTest() {
	suite.NoError(suite.RedisRepo.RedisMasterClient.FlushDB().Err())
}

func (suite *RedisURLRepoSuite) TearDownTest() {
	suite.NoError(suite.RedisRepo.RedisMasterClient.FlushDB().Err())
}

func (suite *RedisURLRepoSuite) TestCreateAndFind() {
	url := model.URL{
		URL: "https://github.com/saraghaedi",
	}

	err := suite.RedisRepo.Create(&url)
	suite.NoError(err)

	id := url.ID

	findByIDUrl, err := suite.RedisRepo.FindByID(uint64(id))
	suite.NoError(err)
	suite.NotNil(findByIDUrl)
	suite.Equal(url.URL, findByIDUrl.URL)

	findByIDUrl, err = suite.RedisRepo.FindByID(uint64(id))
	suite.NoError(err)
	suite.NotNil(findByIDUrl)
	suite.Equal(url.URL, findByIDUrl.URL)
}

func (suite *RedisURLRepoSuite) TestUpdateCounter() {
	url := model.URL{
		URL: "https://github.com/saraghaedi",
	}

	err := suite.RedisRepo.Create(&url)
	suite.NoError(err)

	id := url.ID

	err = suite.RedisRepo.Update(uint64(id), 10)
	suite.NoError(err)

	updatedURL, err := suite.RedisRepo.FindByID(uint64(id))
	suite.NoError(err)

	suite.EqualValues(10, updatedURL.Count)
}

func TestRedisURLRepoSuite(t *testing.T) {
	suite.Run(t, new(RedisURLRepoSuite))
}
