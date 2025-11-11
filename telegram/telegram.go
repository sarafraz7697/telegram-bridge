package telegram

import (
	"encoding/base64"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bridge/logger"
)

// Service handles Telegram bot operations
type Service struct {
	bot    *tgbotapi.BotAPI
	logger logger.Logger
}

// NewService creates a new Telegram service
func NewService(token string, log logger.Logger) (*Service, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	log.Info("Telegram bot authorized", logger.String("username", bot.Self.UserName))

	return &Service{
		bot:    bot,
		logger: log,
	}, nil
}

// SendMessage sends a text message to specified user IDs
func (s *Service) SendMessage(userIDs []int64, message string) error {
	for _, userID := range userIDs {
		msg := tgbotapi.NewMessage(userID, message)
		msg.ParseMode = "HTML"

		_, err := s.bot.Send(msg)
		if err != nil {
			s.logger.Error("Failed to send message",
				logger.Int64("user_id", userID),
				logger.String("error", err.Error()))
			return fmt.Errorf("failed to send message to user %d: %w", userID, err)
		}

		s.logger.Info("Message sent successfully", logger.Int64("user_id", userID))
	}

	return nil
}

// SendImageBase64 sends an image from base64 data to specified user IDs
func (s *Service) SendImageBase64(userIDs []int64, imageBase64 string, caption string) error {
	// Decode base64 image
	imageBytes, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %w", err)
	}

	for _, userID := range userIDs {
		photo := tgbotapi.NewPhoto(userID, tgbotapi.FileBytes{
			Name:  "image.jpg",
			Bytes: imageBytes,
		})
		if caption != "" {
			photo.Caption = caption
			photo.ParseMode = "HTML"
		}

		_, err := s.bot.Send(photo)
		if err != nil {
			s.logger.Error("Failed to send image",
				logger.Int64("user_id", userID),
				logger.String("error", err.Error()))
			return fmt.Errorf("failed to send image to user %d: %w", userID, err)
		}

		s.logger.Info("Image sent successfully",
			logger.Int64("user_id", userID))
	}

	return nil
}

// SendImageWithMessage sends an image with a message to specified user IDs
func (s *Service) SendImageWithMessage(userIDs []int64, imageBase64 string, message string) error {
	return s.SendImageBase64(userIDs, imageBase64, message)
}
