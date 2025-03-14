package redis

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Client 定义了一个Redis客户端结构体
type Client struct {
	*redis.Client
}

var globalClient *Client // 全局变量用于存储Redis客户端实例

// InitRedis 初始化Redis连接
func InitRedis() {

	// globalClient = NewClient("10.192.168.215:6379", "", 0)
	globalClient = NewClient("192.168.123.1:6379", "", 0)
	if err := globalClient.Ping(ctx).Err(); err != nil { // 假设ctx已经定义
		log.Fatalf("cannot connect to redis: %v", err)
	}
	log.Println("success connect Redis")
}

// GetClient 返回全局的Redis客户端实例
func GetClient() *Client {
	return globalClient
}

// NewClient 创建一个新的Redis客户端实例
func NewClient(addr, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{rdb}
}

// Get 从Redis中获取数据
func (c *Client) Get(key string) (string, error) {
	val, err := c.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // 缓存未命中
	} else if err != nil {
		return "", err // 其他错误
	}
	return val, nil
}

// Set 将数据存储到Redis中
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// Select 切换当前使用的数据库
func (c *Client) Select(db int) error {
	_, err := c.Client.Do(ctx, "SELECT", db).Result()
	return err
}
