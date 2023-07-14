package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"srun4-pay/init/common"
	"sync"
	"time"
)

var (
	Rdb16380 *redis.Client
	Rdb16382 *redis.Client
	Rdb16384 *redis.Client
	ctx      = context.Background()
	wg       sync.WaitGroup
)

func init() {
	wg.Add(3)
	go func() {
		defer wg.Done()
		Rdb16380 = createClient(common.Conf.OnlineServer, "16380")
	}()
	go func() {
		defer wg.Done()
		Rdb16382 = createClient(common.Conf.UserServer, "16382")
	}()
	go func() {
		defer wg.Done()
		Rdb16384 = createClient(common.Conf.CacheServer, "16384")
	}()
	wg.Wait()
}

// 创建 Redis 连接
func createClient(host, port string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     common.Conf.RedisPassword,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     100,
		PoolTimeout:  30 * time.Second,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		common.Log.WithField(fmt.Sprintf("Redis[%s] init err", port), err).Warn()
		return nil
	}
	common.Log.WithField(fmt.Sprintf("Redis[:%s] init", port), "Successful").Info()
	return rdb
}
