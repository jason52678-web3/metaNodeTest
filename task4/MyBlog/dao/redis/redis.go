package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	rds  *redis.Client
	once sync.Once
)

func Init() (err error) {
	rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", viper.GetString("redis.host"), viper.GetInt("redis.port")),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = rds.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
	}

	return err

}

func Close() error {
	return rds.Close()
}
