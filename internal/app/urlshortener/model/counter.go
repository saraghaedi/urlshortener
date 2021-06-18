package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// URLCounterRepo represent repository model for counter.
type URLCounterRepo interface {
	Incr(id uint64) error
	Decr(id uint64, count int64) error
	Keys() ([]string, error)
	Get(key string) (int64, error)
}

// RedisURLCounterRepo represent repository model for redis.
type RedisURLCounterRepo struct {
	RedisMasterClient redis.Cmdable
	RedisSlaveClient  redis.Cmdable
}

func (r RedisURLCounterRepo) keyFormat(id uint64) string {
	return fmt.Sprintf("urlshortener:view:%d", id)
}

// Incr increases a counter value each time a shorted URL is called.
func (r RedisURLCounterRepo) Incr(id uint64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(nosqlRepoName, "incr", startTime, finalErr) }()

	if err := r.RedisMasterClient.SetNX(r.keyFormat(id), 0, time.Hour).Err(); err != nil {
		return err
	}

	return r.RedisMasterClient.Incr(r.keyFormat(id)).Err()
}

// Decr decreases the counter value each time a shorted URL is called.
func (r RedisURLCounterRepo) Decr(id uint64, count int64) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(nosqlRepoName, "decr", startTime, finalErr) }()

	return r.RedisMasterClient.DecrBy(r.keyFormat(id), count).Err()
}

// Keys returns all keys from redis.
func (r RedisURLCounterRepo) Keys() (_ []string, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(nosqlRepoName, "keys", startTime, finalErr) }()

	return r.RedisSlaveClient.Keys("urlshortener:view:*").Result()
}

// Get returns value of a key.
func (r RedisURLCounterRepo) Get(key string) (_ int64, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(nosqlRepoName, "get", startTime, finalErr) }()

	result, err := r.RedisSlaveClient.Get(key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(result, 0, 64)
}
