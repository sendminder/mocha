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
func (s *server) GetUserChannels() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// userID가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid user ID"})
		}

		channels, err := s.rdb.GetUserChannels(userID)
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
func (s *server) GetChannel() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		channelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// channelID가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid channel ID"})
		}

		channel, err := s.rdb.GetChannelByID(channelID)
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
func (s *server) CreateChannel() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		var cc types.CreateChannel
		err := c.BodyParser(&cc)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid request payload"})
		}

		var channel = types.Channel{
			Type:            "dm",
			Name:            cc.Name,
			HostUserID:      cc.HostUserID,
			LastMessageID:   0,
			LastDecryptedID: 0,
			CreatedAt:       time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
		}
		createdChannel, err := s.rdb.CreateChannel(&channel)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to create channel"})
		}

		err = s.cache.SetJoinedUsers(createdChannel.ID, cc.JoinedUsers)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to set joined users in cache"})
		}

		// channel user 생성
		for _, value := range cc.JoinedUsers {
			var cuser = types.ChannelUser{
				ChannelID:         createdChannel.ID,
				UserID:            value,
				LastSeenMessageID: 0,
			}
			err = s.rdb.CreateChannelUser(&cuser)
			if err != nil {
				return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to create channel_user"})
			}
		}
		return c.JSON(map[string]types.Channel{"channel": *createdChannel})
	}
}
