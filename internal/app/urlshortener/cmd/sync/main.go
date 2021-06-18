package sync

import (
	"strconv"
	"strings"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/metric"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	syncInterval        = 5 * time.Second
	healthCheckInterval = 1
)

func main(cfg config.Config) {
	masterDb, err := database.New(cfg.Database.Driver, cfg.Database.MasterConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to master database: %s", err.Error())
	}

	slaveDb, err := database.New(cfg.Database.Driver, cfg.Database.SlaveConnStr)
	if err != nil {
		logrus.Fatalf("faled to connect to slave database: %s", err.Error())
	}

	redisCfg := cfg.Redis

	redisMasterClient, redisMasterClose := redis.Create(redisCfg.MasterAddress, redisCfg.Options, true)
	redisSlaveClient, redisSlaveClose := redis.Create(redisCfg.SlaveAddress, redisCfg.Options, false)

	defer func() {
		if err := redisMasterClose(); err != nil {
			logrus.Errorf("redis master connection close error: %s", err.Error())
		}

		if err := redisSlaveClose(); err != nil {
			logrus.Errorf("redis slave connection close error: %s", err.Error())
		}
	}()

	_, err1 := scheduler.Every(healthCheckInterval).Seconds().Run(func() {
		metric.ReportDbStatus(masterDb, "database_master")
		metric.ReportDbStatus(slaveDb, "database_slave")
		metric.ReportRedisStatus(redisMasterClient, "redis_master")
		metric.ReportRedisStatus(redisSlaveClient, "redis_slave")
	})
	if err1 != nil {
		logrus.Fatalf("failed to start metric scheduler: %s", err1.Error())
	}

	urlRepo := model.SQLURLRepo{
		MasterDB: masterDb,
		SlaveDB:  slaveDb,
	}

	urlCounterRepo := model.RedisURLCounterRepo{
		RedisMasterClient: redisMasterClient,
		RedisSlaveClient:  redisSlaveClient,
	}

	scheduleWorker(urlRepo, urlCounterRepo)
}

func scheduleWorker(urlRepo model.URLRepo, urlCounterRepo model.URLCounterRepo) {
	for {
		time.Sleep(syncInterval)

		if err := runWorker(urlRepo, urlCounterRepo); err != nil {
			logrus.Errorf("failed to run worker: %s", err.Error())
		}
	}
}

func runWorker(urlRepo model.URLRepo, urlCounterRepo model.URLCounterRepo) error {
	keys, err := urlCounterRepo.Keys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := updateCounter(key, urlRepo, urlCounterRepo); err != nil {
			logrus.Errorf("failed to update url view count: %s", err.Error())
		}
	}

	return nil
}

func updateCounter(key string, urlRepo model.URLRepo, urlCounterRepo model.URLCounterRepo) error {
	count, err := urlCounterRepo.Get(key)
	if err != nil {
		return err
	}

	if count == 0 {
		return nil
	}

	idStr := strings.Split(key, ":")[2]

	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		return err
	}

	if err := urlRepo.Update(uint64(id), count); err != nil {
		return err
	}

	if err := urlCounterRepo.Decr(uint64(id), count); err != nil {
		return err
	}

	return nil
}

// Register registers sync command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync database with Redis url counts",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}
	root.AddCommand(cmd)
}
