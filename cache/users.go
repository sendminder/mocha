package cache

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func init() {
	// Redis 클라이언트 초기화
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:26379", // Redis 서버 주소
		Password: "",                // 비밀번호 (없는 경우 빈 문자열)
		DB:       0,                 // 데이터베이스 번호
	})

	// 연결 확인
	pong, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Connected to Redis:", pong)
}

func SetJoinedUsers(convId int64, userIds []int64) error {
	// userIds를 문자열로 변환하여 Redis에 저장
	userIdsStr := make([]string, len(userIds))
	for i, id := range userIds {
		userIdsStr[i] = fmt.Sprintf("%d", id)
	}

	key := fmt.Sprintf("conversation:%d:users", convId)
	err := RedisClient.SAdd(context.Background(), key, userIdsStr).Err()
	if err != nil {
		return err
	}
	log.Println("[Redis]SetJoinedUsers:", userIds)
	return nil
}

func GetJoinedUsers(convId int64) ([]int64, error) {
	key := fmt.Sprintf("conversation:%d:users", convId)
	userIdsStr, err := RedisClient.SMembers(context.Background(), key).Result()
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
