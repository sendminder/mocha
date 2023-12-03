package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type ChannelHandler interface {
	GetUserChannels() func(*fiber.Ctx) error
	GetChannel() func(*fiber.Ctx) error
	CreateChannel() func(*fiber.Ctx) error
}

// GetChannelsHandler는 해당 유저의 모든 채팅방을 반환합니다.
func (s *restServer) GetUserChannels() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		userId, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// userId가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid user Id"})
		}

		channels, err := s.rdb.GetUserChannels(userId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 레코드를 찾지 못한 경우 404 에러 반환
				return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "Channel not found"})
			}
			// 다른 에러가 발생한 경우 500 에러 반환
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to get channel"})
		}

		return c.JSON(map[string][]types.Channel{"channels": channels})
	}
}

// GetChannelHandler는 특정 채팅방을 반환합니다.
func (s *restServer) GetChannel() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		channelId, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// channelId가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid channel Id"})
		}

		channel, err := s.rdb.GetChannelByID(channelId)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 레코드를 찾지 못한 경우 404 에러 반환
				return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "Channel not found"})
			}
			// 다른 에러가 발생한 경우 500 에러 반환
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to get channel"})
		}
		return c.JSON(map[string]types.Channel{"channel": *channel})
	}
}

// CreateChannelHandler는 새로운 채팅방을 생성합니다.
func (s *restServer) CreateChannel() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		var cc types.CreateChannel
		_ = c.BodyParser(&cc)

		var channel = types.Channel{
			Type:            "dm",
			Name:            cc.Name,
			HostUserId:      cc.HostUserId,
			LastMessageId:   0,
			LastDecryptedId: 0,
			CreatedAt:       time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
		}
		createdChannel, err := s.rdb.CreateChannel(&channel)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to create channel"})
		}

		err = s.cache.SetJoinedUsers(createdChannel.Id, cc.JoinedUsers)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to set joined users in cache"})
		}

		// channel user 생성
		for _, value := range cc.JoinedUsers {
			var cuser = types.ChannelUser{
				ChannelId:         createdChannel.Id,
				UserId:            value,
				LastSeenMessageId: 0,
			}
			err = s.rdb.CreateChannelUser(&cuser)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to create channel_user"})
			}
		}
		return c.JSON(map[string]types.Channel{"channel": *createdChannel})
	}
}
