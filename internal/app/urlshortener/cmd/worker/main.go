package worker

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/carlescere/scheduler"
	"github.com/nats-io/nats.go"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/metric"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/database"
	pkgnats "github.com/saraghaedi/urlshortener/pkg/nats"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	healthCheckInterval = 1
)

// nolint:funlen
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

	natsConn, err := pkgnats.Create(cfg.Nats)
	if err != nil {
		logrus.Fatalf("Failed to create nats connection: %s", err.Error())
	}

	defer natsConn.Close()

	_, err = natsConn.QueueSubscribe(config.App, config.App, func(msg *nats.Msg) {
		if err := updateCounter(string(msg.Data), urlRepo, urlCounterRepo); err != nil {
			logrus.Errorf("failed to update url counter: %s", err)
		}
	})
	if err != nil {
		logrus.Fatalf("failed to subscribe to nats: %s", err.Error())
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig

	logrus.Infof("signal %s received", s)
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

// Register registers worker command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "worker",
		Short: "Update URL counter for each URL",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}
	root.AddCommand(cmd)
}
