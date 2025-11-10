package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"telegram-bridge/logger"
	"telegram-bridge/models"
)

// Client handles communication with Telegram Bot API
type Client struct {
	httpClient *http.Client
	logger     logger.Logger
}

// NewClient creates a new Telegram client
func NewClient(log logger.Logger) *Client {
	return &Client{
		httpClient: &http.Client{},
		logger:     log,
	}
}

// SendMessage sends a message to a specific chat via Telegram Bot API
func (c *Client) SendMessage(botToken, chatID, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	payload := models.TelegramMessage{
		ChatID: chatID,
		Text:   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		c.logger.Error("Failed to marshal telegram payload",
			logger.ZError(err),
			logger.String("chat_id", chatID))
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to send request to Telegram API",
			logger.ZError(err),
			logger.String("chat_id", chatID))
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Telegram API returned error",
			logger.Int("status_code", resp.StatusCode),
			logger.String("response", string(body)),
			logger.String("chat_id", chatID))
		return fmt.Errorf("telegram API returned status %d: %s", resp.StatusCode, string(body))
	}

	c.logger.Info("Message sent successfully",
		logger.String("chat_id", chatID))
	return nil
}

// SendMessageToMultiple sends a message to multiple users concurrently
func (c *Client) SendMessageToMultiple(botToken, message string, userIDs []string) (successCount, failCount int) {
	var wg sync.WaitGroup
	var mu sync.Mutex

	c.logger.Info("Starting batch message send",
		logger.Int("recipient_count", len(userIDs)))

	for _, userID := range userIDs {
		wg.Add(1)
		go func(uid string) {
			defer wg.Done()
			if err := c.SendMessage(botToken, uid, message); err != nil {
				mu.Lock()
				failCount++
				mu.Unlock()
			} else {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(userID)
	}

	wg.Wait()

	c.logger.Info("Batch message send completed",
		logger.Int("success_count", successCount),
		logger.Int("fail_count", failCount))

	return successCount, failCount
}
