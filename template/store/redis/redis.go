package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"survey-ussd/core/environment"
	"time"
)

type Cache struct {
	env    *environment.Environment
	ctx    context.Context
	client *redis.Client
}

func New(env *environment.Environment, ctx context.Context) *Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     env.RedisURL,
		Password: env.RedisPassword,
		DB:       0, // use default DB
	})
	return &Cache{env: env, ctx: ctx, client: rdb}
}

func (conf *Cache) SetValue(key string, value interface{}, duration time.Duration) error {
	err := conf.client.Set(conf.ctx, key, value, duration).Err()
	//err := conf.client.Set(conf.ctx, key, value, 5*time.Hour).Err()
	if err != nil {
		log.Fatal(err)
	}
	return err
}

func (conf *Cache) GetValue(key string) (string, error) {
	val, err := conf.client.Get(conf.ctx, key).Result()
	return val, err
}

func (conf *Cache) GetKeys(pattern string) ([]string, error) {
	var keys []string
	iter := conf.client.Scan(conf.ctx, 0, pattern, 0).Iterator()
	for iter.Next(conf.ctx) {
		keys = append(keys, iter.Val())
	}
	return keys, nil
}
func (conf *Cache) DeleteValue(key string) error {
	err := conf.client.Del(conf.ctx, key).Err()
	return err
}
