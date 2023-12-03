package cache

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type RedisCache interface {
	SetJoinedUsers(channelID int64, userIDs []int64) error
	GetJoinedUsers(channelID int64) ([]int64, error)
}

var _ RedisCache = (*redisCache)(nil)

type redisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache(ctx context.Context, host string) RedisCache {
	// Redis 클라이언트 초기화
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host, // Redis 서버 주소
		Password: "",   // 비밀번호 (없는 경우 빈 문자열)
		DB:       0,    // 데이터베이스 번호
	})

	// localhost:26379
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		return nil
	}
	slog.Info("Connected to Redis", "pong", pong)
	return &redisCache{
		client: redisClient,
		ctx:    ctx,
	}
}

func (r *redisCache) SetJoinedUsers(channelID int64, userIDs []int64) error {
	// userIDs를 문자열로 변환하여 Redis에 저장
	userIDsStr := make([]string, len(userIDs))
	for i, id := range userIDs {
		userIDsStr[i] = fmt.Sprintf("%d", id)
	}

	key := fmt.Sprintf("channel:%d:users", channelID)
	err := r.client.SAdd(r.ctx, key, userIDsStr).Err()
	if err != nil {
		return err
	}
	slog.Info("[Redis]SetJoinedUsers", "userIDs", userIDs)
	return nil
}

func (r *redisCache) GetJoinedUsers(channelID int64) ([]int64, error) {
	key := fmt.Sprintf("channel:%d:users", channelID)
	userIDsStr, err := r.client.SMembers(r.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	userIDs := make([]int64, len(userIDsStr))
	for i, idStr := range userIDsStr {
		var id int64
		_, err := fmt.Sscanf(idStr, "%d", &id)
		if err != nil {
			return nil, err
		}
		userIDs[i] = id
	}
	slog.Info("[Redis]GetJoinedUsers:", "userIDs", userIDs)
	return userIDs, nil
}
