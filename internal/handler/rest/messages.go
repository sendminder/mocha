package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type MessageRestHandler interface {
	GetMessages() func(*fiber.Ctx) error
}

func (s *server) GetMessages() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		channelID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// channelID가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid channel ID"})
		}
		limit, err := strconv.ParseInt(c.Params("limit"), 10, 64)
		if err != nil {
			// limit이 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid limit"})
		}

		messages, err := s.mdb.GetMessagesByChannelID(channelID, limit)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 레코드를 찾지 못한 경우 404 에러 반환
				return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "messages not found"})
			}
			// 다른 에러가 발생한 경우 500 에러 반환
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to get messages"})
		}
		return c.JSON(map[string][]types.Message{"messages": messages})
	}
}
