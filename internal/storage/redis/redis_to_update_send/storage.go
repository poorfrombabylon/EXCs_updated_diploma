package redis_to_update_send

import (
	"context"
	"encoding/json"
	"excs_updater/internal/config"
	redisQueueDomain "excs_updater/internal/domain/queues"
	"fmt"
	"github.com/go-redis/redis"
)

type RedisToUpdateSend interface {
	SendLinkedinID(_ context.Context, profile redisQueueDomain.ProfileToUpdate) error
}

type queueConverter struct {
	client         *redis.Client
	topic          string
	deliveryMethod string
}

func NewRedisToUpdateSend(cfg config.RedisConfig) (RedisToUpdateSend, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", cfg.Host, cfg.Port),
		Password: cfg.Password,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	return &queueConverter{
		client:         client,
		topic:          cfg.Topic,
		deliveryMethod: cfg.DeliveryMethod,
	}, nil
}

func (c *queueConverter) SendLinkedinID(_ context.Context, profile redisQueueDomain.ProfileToUpdate) error {
	body, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	err = c.client.SAdd(c.topic, body).Err()

	return err
}
