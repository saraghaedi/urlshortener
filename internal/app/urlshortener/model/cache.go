package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const expirationTimeInterval = time.Hour

// RedisURLRepo represent repository model for redis cache.
type RedisURLRepo struct {
	Base              URLRepo
	RedisMasterClient redis.Cmdable
	RedisSlaveClient  redis.Cmdable
}

func (r RedisURLRepo) format(id uint64) string {
	return fmt.Sprintf("urlshortener:url:id:%d", id)
}

// Create creates a new shorted url in database.
func (r RedisURLRepo) Create(url *URL) error {
	return r.Base.Create(url)
}

// FindByID finds a url in cache by ID.
func (r RedisURLRepo) FindByID(id uint64) (*URL, error) {
	result, err := r.RedisSlaveClient.Get(r.format(id)).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	} else if err == nil {
		url := URL{}

		err = json.Unmarshal([]byte(result), &url)
		if err != nil {
			return nil, err
		}

		return &url, nil
	}

	value, err := r.Base.FindByID(id)
	if err != nil {
		return nil, err
	}

	jsonVal, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	err = r.RedisMasterClient.Set(r.format(id), jsonVal, expirationTimeInterval).Err()
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Update updates count column in database.
func (r RedisURLRepo) Update(id uint64, additionalCount int64) error {
	return r.Base.Update(id, additionalCount)
}
