package models

// SendMessageRequest represents the incoming API request
type SendMessageRequest struct {
	BotToken string   `json:"bot_token"`
	Message  string   `json:"message"`
	ChatID   string   `json:"chat_id,omitempty"`  // Optional: specific chat_id to send to
	UserIDs  []string `json:"user_ids,omitempty"` // Optional: list of user IDs to send to
}

// TelegramMessage represents the message to send to Telegram API
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// SendMessageResponse represents the response for send message requests
type SendMessageResponse struct {
	Success      bool `json:"success"`
	TotalSent    int  `json:"total_sent"`
	SuccessCount int  `json:"success_count"`
	FailCount    int  `json:"fail_count"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}
