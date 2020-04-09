package redis

import (
	"log"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	InitRedisDB()
}

// InitRedisDB 实现RedisClient的实例化，单例模式，饿汉模式
func InitRedisDB() {
	// client是多线程安全的，可以在多个goroutine中使用，不需要额外加锁
	if client == nil {
		log.Println("Redis client firstly initializing...")
		client = redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379",
			Password: "",
			DB:       0,
		})
		log.Println("Redis client firstly initialized")
	}

	pong, err := client.Ping().Result()
	log.Println("Pint test: ...", pong, err)
}

// HKExists 测试key中的field是否存在
func HKExists(key string, field string) (bool, error) {
	if client == nil {
		InitRedisDB()
	}
	return client.HExists(key, field).Result()
}

// HIncr 将key中指定field的value增加1
func HIncr(key string, field string) (int64, error) {
	if client == nil {
		InitRedisDB()
	}
	return client.HIncrBy(key, field, 1).Result()
}

// HSetNX 用来插入没有的field，如果有则返回错误
func HSetNX(key string, field string) (bool, error) {
	if client == nil {
		InitRedisDB()
	}
	return client.HSetNX(key, field, 1).Result()
}

// Close used for close the client
func Close() {
	if client != nil {
		client.Close()
	}
}
