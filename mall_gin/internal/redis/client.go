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
type MemberScore struct {
	Member string
	Score  float64
}

var globalClient *Client // 全局变量用于存储Redis客户端实例

// InitRedis 初始化Redis连接
func InitRedis() {
	redishost := "192.168.123.1:6379"
	globalClient = NewClient(redishost, "", 0)
	if err := globalClient.Ping(ctx).Err(); err != nil {
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

// []MemberScore
func (c *Client) GetFromZSET(key string, n int) ([]MemberScore, error) {
	// 前n个成员
	// 使用ZRANGE获取分数最低的前n个成员及其分数
	membersWithScores, err := c.Client.ZRevRangeWithScores(ctx, key, 0, int64(n-1)).Result()
	if err != nil {
		log.Fatalf("Failed to retrieve members: %v", err)
		return nil, err
	}
	// return membersWithScores, nil
	var result []MemberScore
	for _, member := range membersWithScores {
		// 确保Member字段是字符串类型
		memberStr, ok := member.Member.(string)
		if !ok {
			log.Printf("Member is not a string: %v", member.Member)
			continue
		}
		result = append(result, MemberScore{
			Member: memberStr,
			Score:  member.Score,
		})
	}

	return result, nil
}

// Set 将数据存储到Redis中
func (c *Client) Set(key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}

// ZSet 将数据存储到Redis中并更新
func (c *Client) ZSetUpgrade(key string, memberName string, expiration time.Duration) error {
	// 初始化上下文
	_, err := c.Client.ZScore(ctx, key, memberName).Result()
	if err == redis.Nil {
		// 成员不存在，添加新成员并设置分数为1
		err := c.Client.ZAdd(ctx, key, &redis.Z{
			Score:  1,
			Member: memberName,
		}).Err()
		if err != nil {
			log.Fatalf("Failed to add new member: %v", err)
			return err
		}
		log.Println("New member added with score 1")
		return nil
	} else if err != nil {
		log.Fatalf("Failed to check member existence: %v", err)
		return err
	} else {
		// 成员存在，分数自增1
		newScore, err := c.Client.ZIncrBy(ctx, key, 1, memberName).Result()
		if err != nil {
			log.Fatalf("Failed to increment score: %v", err)
			return err
		}
		log.Printf("Member's score incremented to %f\n", newScore)
		return nil
	}
}

// Select 切换当前使用的数据库
func (c *Client) Select(db int) error {
	_, err := c.Client.Do(ctx, "SELECT", db).Result()
	return err
}
