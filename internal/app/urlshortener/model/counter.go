package model

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

// URLCounterRepo represent repository model for counter.
type URLCounterRepo interface {
	Incr(id uint64) error
}

// RedisURLCounterRepo represent repository model for redis.
type RedisURLCounterRepo struct {
	RedisMasterClient redis.Cmdable
	RedisSlaveClient  redis.Cmdable
}

// Incr will increase a counter value each time a shorted URL is called.
func (r RedisURLCounterRepo) Incr(id uint64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(nosqlRepoName, "Counter", startTime, finalErr) }()

	return r.RedisMasterClient.Incr(fmt.Sprintf("%d", id)).Err()
}
