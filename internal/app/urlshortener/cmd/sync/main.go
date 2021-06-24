package sync

import (
	"time"

	pkgnats "github.com/saraghaedi/urlshortener/pkg/nats"

	"github.com/carlescere/scheduler"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/metric"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/model"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/nats-io/nats.go"
)

const (
	syncInterval        = 5 * time.Minute
	healthCheckInterval = 1
)

func main(cfg config.Config) {
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
		metric.ReportRedisStatus(redisMasterClient, "redis_master")
		metric.ReportRedisStatus(redisSlaveClient, "redis_slave")
	})
	if err1 != nil {
		logrus.Fatalf("failed to start metric scheduler: %s", err1.Error())
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

	scheduleWorker(urlCounterRepo, natsConn)
}

func scheduleWorker(urlCounterRepo model.URLCounterRepo, natsConn *nats.Conn) {
	for {
		time.Sleep(syncInterval)

		if err := runWorker(urlCounterRepo, natsConn); err != nil {
			logrus.Errorf("failed to run worker: %s", err.Error())
		}
	}
}

func runWorker(urlCounterRepo model.URLCounterRepo, natsConn *nats.Conn) error {
	keys, err := urlCounterRepo.Keys()
	if err != nil {
		return err
	}

	for _, key := range keys {
		if err := natsConn.Publish(config.App, []byte(key)); err != nil {
			logrus.Errorf("failed to publish key to nats: %s", err.Error())
		}
	}

	return nil
}

// Register registers sync command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Schedule URL keys for updating its counter",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}
	root.AddCommand(cmd)
}
