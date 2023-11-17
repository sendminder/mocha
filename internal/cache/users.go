package cache

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	SetJoinedUsers(convId int64, userIds []int64) error
	GetJoinedUsers(convId int64) ([]int64, error)
}

var _ RedisCache = (*redisCache)(nil)

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(host string) RedisCache {
	// Redis 클라이언트 초기화
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host, // Redis 서버 주소
		Password: "",   // 비밀번호 (없는 경우 빈 문자열)
		DB:       0,    // 데이터베이스 번호
	})

	// localhost:26379
	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		slog.Error("Failed to connect to Redis:", "error", err)
	}
	slog.Info("Connected to Redis", "pong", pong)
	return &redisCache{
		client: redisClient,
	}
}

func (r *redisCache) SetJoinedUsers(convId int64, userIds []int64) error {
	// userIds를 문자열로 변환하여 Redis에 저장
	userIdsStr := make([]string, len(userIds))
	for i, id := range userIds {
		userIdsStr[i] = fmt.Sprintf("%d", id)
	}

	key := fmt.Sprintf("conversation:%d:users", convId)
	err := r.client.SAdd(context.Background(), key, userIdsStr).Err()
	if err != nil {
		return err
	}
	log.Println("[Redis]SetJoinedUsers:", userIds)
	return nil
}

func (r *redisCache) GetJoinedUsers(convId int64) ([]int64, error) {
	key := fmt.Sprintf("conversation:%d:users", convId)
	userIdsStr, err := r.client.SMembers(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	userIds := make([]int64, len(userIdsStr))
	for i, idStr := range userIdsStr {
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			return nil, err
		}
		userIds[i] = id
	}
	log.Println("[Redis]GetJoinedUsers:", userIds)
	return userIds, nil
}
