package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"mocha/internal/types"
)

type UserHandler interface {
	GetUser() func(*fiber.Ctx) error
	CreateUser() func(*fiber.Ctx) error
	LoginUser() func(*fiber.Ctx) error
}

func (s *server) GetUser() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		userID, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			// userID가 올바른 int64로 변환되지 않은 경우 에러 처리
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid user ID"})
		}

		user, err := s.rdb.GetUser(userID)
		if err != nil {
			return s.handleError(c, err)
		}
		return c.JSON(map[string]types.User{"user": *user})
	}
}

func (s *server) CreateUser() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		var cu types.CreateUser
		err := c.BodyParser(&cu)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid request payload"})
		}

		// validator를 사용하여 필수 파라미터 체크
		validate := validator.New()
		if err := validate.Struct(cu); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		foundUser, err := s.rdb.GetUserByEmail(cu.Email)
		if err != nil {
			return s.handleError(c, err)
		}
		if foundUser != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Duplicated user"})
		}

		var user = types.User{
			Name:      cu.Name,
			Password:  cu.Password,
			Email:     cu.Email,
			Age:       cu.Age,
			Gender:    cu.Gender,
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
			UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		}

		createdUser, err := s.rdb.CreateUser(&user)
		if err != nil {
			return s.handleError(c, err)
		}

		// create bot chat
		/*
			1. bot user 가져오기 이름 meow
			2. bot user의 user id로 chat방 만들기
		*/

		botUser, err := s.rdb.GetBotByName("meow")
		if err != nil {
			return s.handleError(c, err)
		}
		var channel = types.Channel{
			Type:            "bot",
			Name:            "meow-meow",
			HostUserID:      user.ID,
			LastMessageID:   0,
			LastDecryptedID: 0,
			CreatedAt:       time.Now().UTC().Format(time.RFC3339),
			UpdatedAt:       time.Now().UTC().Format(time.RFC3339),
		}
		createdBotChannel, err := s.rdb.CreateChannel(&channel)
		if err != nil {
			return s.handleError(c, err)
		}

		// channel user 생성
		var cuser = types.ChannelUser{
			ChannelID:         createdBotChannel.ID,
			UserID:            user.ID,
			LastSeenMessageID: 0,
		}
		err = s.rdb.CreateChannelUser(&cuser)
		if err != nil {
			return s.handleError(c, err)
		}

		cuser = types.ChannelUser{
			ChannelID:         createdBotChannel.ID,
			UserID:            botUser.ID,
			LastSeenMessageID: 0,
		}
		err = s.rdb.CreateChannelUser(&cuser)
		if err != nil {
			return s.handleError(c, err)
		}
		return c.JSON(map[string]types.User{"user": *createdUser})
	}
}

func (s *server) LoginUser() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		var lu types.LoginUser
		err := c.BodyParser(&lu)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid request payload"})
		}

		// validator를 사용하여 필수 파라미터 체크
		validate := validator.New()
		if err := validate.Struct(lu); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": err.Error()})
		}

		foundUser, err := s.rdb.GetUserByEmail(lu.Email)
		if err != nil {
			return s.handleError(c, err)
		}
		if foundUser == nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Duplicated user"})
		}

		if foundUser.Password != lu.Password {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{"error": "Invalid Password"})
		}
		return c.JSON(map[string]types.User{"user": *foundUser})
	}
}

func (s *server) handleError(c *fiber.Ctx, err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 레코드를 찾지 못한 경우 404 에러 반환
		return c.Status(http.StatusNotFound).JSON(map[string]string{"error": "User not found"})
	}
	// 다른 에러가 발생한 경우 500 에러 반환
	return c.Status(http.StatusInternalServerError).JSON(map[string]string{"error": "Failed to get user"})
}
