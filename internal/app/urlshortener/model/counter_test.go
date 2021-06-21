package model_test

import (
	"fmt"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/stretchr/testify/suite"
	"testing"
)

type URLCounterRepoSuite struct {
	suite.Suite
	CounterRepo model.RedisURLCounterRepo
}

func (suite *URLCounterRepoSuite) SetupSuite() {
	cfg := config.Init()

	redisCfg := cfg.Redis

	rMaster, _ := redis.Create(redisCfg.MasterAddress, redisCfg.Options, true)
	rSlave, _ := redis.Create(redisCfg.SlaveAddress, redisCfg.Options, false)

	suite.CounterRepo = model.RedisURLCounterRepo{
		RedisMasterClient: rMaster,
		RedisSlaveClient:  rSlave,
	}
}

func (suite *URLCounterRepoSuite) SetupTest() {
	suite.NoError(suite.CounterRepo.RedisMasterClient.FlushDB().Err())
}

func (suite *URLCounterRepoSuite) TearDownTest() {
	suite.NoError(suite.CounterRepo.RedisMasterClient.FlushDB().Err())
}

func (suite *URLCounterRepoSuite) TestURLCounterRepo() {
	var id uint64 = 1

	key := fmt.Sprintf("urlshortener:view:%d", id)

	err := suite.CounterRepo.Incr(id)
	suite.NoError(err)

	value, err := suite.CounterRepo.Get(key)
	suite.EqualValues(1, value)

	err = suite.CounterRepo.Decr(id, 1)
	suite.NoError(err)

	value, err = suite.CounterRepo.Get(key)
	suite.EqualValues(0, value)

	keys, err := suite.CounterRepo.Keys()
	keysList := []string{key}
	suite.Equal(keysList, keys)
}

func TestURLCounterRepoSuite(t *testing.T) {
	suite.Run(t, new(URLCounterRepoSuite))
}
