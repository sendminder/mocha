package rest

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"mocha/internal/types"
)

type DeviceHandler interface {
	GetDevice() func(*fiber.Ctx) error
	CreateDevice() func(*fiber.Ctx) error
}

func (s *server) GetDevice() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		deviceID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// deviceID가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid device ID"})
		}

		device, err := s.rdb.GetDevice(deviceID)
		if err != nil {
			return s.handleError(c, err)
		}
		return c.JSON(map[string]types.Device{"device": *device})
	}
}

func (s *server) CreateDevice() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		var cd types.CreateDevice
		err := c.BodyParser(&cd)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid request payload"})
		}

		// validator를 사용하여 필수 파라미터 체크
		validate := validator.New()
		if err := validate.Struct(cd); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		foundDevice, err := s.rdb.GetDeviceByPushToken(cd.PushToken)
		if err != nil {
			return s.handleError(c, err)
		}
		if foundDevice != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Duplicated device"})
		}

		var device = types.Device{
			UserID:    cd.UserID,
			PushToken: cd.PushToken,
			Platform:  cd.Platform,
			Version:   cd.Version,
			Activated: true,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		createdDevice, err := s.rdb.CreateDevice(&device)
		if err != nil {
			return s.handleError(c, err)
		}
		return c.JSON(map[string]types.Device{"device": *createdDevice})
	}
}
