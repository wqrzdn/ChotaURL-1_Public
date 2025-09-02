package storages

import (
	"context"
	"fmt"
	"time"
	"crypto/rand"
	"math/big"
	"github.com/redis/go-redis/v9"
)



type Redis struct {
	client *redis.Client
}


func NewRedis(dsn string) (*Redis, error) {
	
	opt, err := redis.ParseURL(dsn)
	if err != nil {
		return nil, err
	}

	
	rdb := redis.NewClient(opt)

	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Redis{client: rdb}, nil
}


func (r *Redis) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.client.Set(ctx, key, value, expiration).Err()
}


func (r *Redis) Get(key string) (string, error) {
	ctx := context.Background()
	return r.client.Get(ctx, key).Result()
}


func (r *Redis) Code() string {
	return "redis"
}


func (r *Redis) Save(url string) string {
	code := generateCode(url)
	_ = r.Set(code, url, 24*time.Hour)
	return code
}

func (r *Redis) SaveWithCustom(url, customName string) (string, error) {
	
	if len(customName) < 3 || len(customName) > 20 {
		return "", fmt.Errorf("custom name must be 3-20 characters")
	}

	
	if r.Exists(customName) {
		return "", fmt.Errorf("custom name already taken")
	}

	
	err := r.Set(customName, url, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return customName, nil
}

func (r *Redis) Exists(code string) bool {
	ctx := context.Background()
	exists, err := r.client.Exists(ctx, code).Result()
	return err == nil && exists > 0
}


func (r *Redis) Load(code string) (string, error) {
	return r.Get(code)
}


func generateCode(url string) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 6
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			b[i] = letters[0]
		} else {
			b[i] = letters[n.Int64()]
		}
	}
	return string(b)
}
