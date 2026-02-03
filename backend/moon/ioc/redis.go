package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func InitRedis() *redis.Client {
	c := RedisConfig{
		Addr: "localhost:6379",
		DB:   0,
	}

	err := viper.UnmarshalKey("redis", &c)
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password,
		DB:       c.DB,
	})
	return client
}
