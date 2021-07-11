package limit

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

func NewRedisCounter(redisURL string) *RedisCounter {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		log.WithError(err).Fatalf("invalid redis url")
	}

	return &RedisCounter{
		db: redis.NewClient(options),
	}
}

type RedisCounter struct {
	db *redis.Client
}

func (rc *RedisCounter) GetCount(ctx context.Context, ip string) (int, error) {
	count, err := rc.db.Get(ctx, ip).Int()
	if err != nil && err != redis.Nil {
		log.WithField("ip", ip).WithError(err).Error("could not get count")
		return 0, err
	}

	log.WithField("ip", ip).Debugf("got count: %v", count)
	return count, nil
}

func (rc *RedisCounter) Increment(ctx context.Context, ip string, ttl int) error {
	newCount, err := rc.db.Incr(ctx, ip).Result()
	if err != nil {
		log.WithField("ip", ip).WithError(err).Error("could not increment")
		return err
	}

	log.WithField("ip", ip).Debugf("incremented to: %v", newCount)

	if ttl > 0 {
		go rc.scheduleExpiration(ip, ttl)
	}

	return nil
}

func (rc *RedisCounter) scheduleExpiration(ip string, ttl int) {
	expiration := time.Duration(ttl) * time.Second
	err := rc.db.Expire(context.Background(), ip, expiration).Err()
	if err != nil {
		log.WithField("ip", ip).WithError(err).Error("could not expire")
	}
}
