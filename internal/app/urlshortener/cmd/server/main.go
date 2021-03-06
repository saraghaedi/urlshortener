package server

import (
	"github.com/carlescere/scheduler"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/config"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/metric"
	"github.com/saraghaedi/urlshortener/internal/app/urlshortener/router"
	"github.com/saraghaedi/urlshortener/pkg/database"
	"github.com/saraghaedi/urlshortener/pkg/prometheus"
	"github.com/saraghaedi/urlshortener/pkg/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
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

	go prometheus.StartServer(cfg.Monitoring.Prometheus)

	r := router.New(cfg, masterDb, slaveDb, redisMasterClient, redisSlaveClient)
	logrus.Fatal(r.Start(cfg.Server.Address))
}

// Register registers server command for urlshortener binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run url shortener server component",
		Run: func(cmd *cobra.Command, args []string) {
			main(cfg)
		},
	}
	root.AddCommand(cmd)
}
