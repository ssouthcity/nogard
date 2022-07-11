package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/ssouthcity/nogard"
)

type DragonCache struct {
	cache  *redis.Client
	srv    nogard.DragonEncyclopedia
	logger logrus.FieldLogger
}

func NewDragonCache(addr string, srv nogard.DragonEncyclopedia, logger logrus.FieldLogger) *DragonCache {
	client := redis.NewClient(&redis.Options{Addr: addr})

	if err := client.Ping(context.Background()).Err(); err != nil {
		logger.WithError(err).WithField("address", addr).Warn("unable to connect to redis, opting out of cache")
	}

	return &DragonCache{client, srv, logger}
}

func (c *DragonCache) DragonNames() ([]string, error) {
	ctx := context.Background()
	key := "dragons"

	buf, err := c.cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		names, err := c.srv.DragonNames()
		if err != nil {
			return nil, err
		}

		s, err := json.Marshal(names)
		if err != nil {
			return nil, err
		}

		if err := c.cache.Set(ctx, key, s, 24*time.Hour).Err(); err != nil {
			return nil, err
		}

		return names, nil
	}

	var names []string

	if err := json.Unmarshal(buf, &names); err != nil {
		return nil, err
	}

	return names, nil
}

func (c *DragonCache) Dragon(name string) (*nogard.Dragon, error) {
	ctx := context.Background()
	key := fmt.Sprintf("dragons:%s", name)

	buf, err := c.cache.Get(ctx, key).Bytes()
	if err == redis.Nil {
		dragon, err := c.srv.Dragon(name)
		if err != nil {
			return nil, err
		}

		s, err := json.Marshal(dragon)
		if err != nil {
			return nil, err
		}

		if err := c.cache.Set(ctx, key, s, 24*time.Hour).Err(); err != nil {
			return nil, err
		}

		return dragon, nil
	}

	var dragon nogard.Dragon

	if err := json.Unmarshal(buf, &dragon); err != nil {
		return nil, err
	}

	return &dragon, nil
}
