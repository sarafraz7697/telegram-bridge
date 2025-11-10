package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"telegram-bridge/logger"
	"telegram-bridge/models"
	"telegram-bridge/telegram"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	telegramClient *telegram.Client
	logger         logger.Logger
}

// NewHandler creates a new Handler instance
func NewHandler(telegramClient *telegram.Client, log logger.Logger) *Handler {
	return &Handler{
		telegramClient: telegramClient,
		logger:         log,
	}
}

// SendMessage handles POST requests to send messages
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req models.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request body", logger.ZError(err))
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.BotToken == "" {
		h.logger.Warn("Missing bot_token in request")
		http.Error(w, "bot_token is required", http.StatusBadRequest)
		return
	}

	if req.Message == "" {
		h.logger.Warn("Missing message in request")
		http.Error(w, "message is required", http.StatusBadRequest)
		return
	}

	var targetUsers []string

	// Determine target users
	if len(req.UserIDs) > 0 {
		// Send to specific user IDs
		targetUsers = req.UserIDs
		h.logger.Info("Sending to specific users",
			logger.Int("user_count", len(targetUsers)))
	} else if req.ChatID != "" {
		// Send to single chat ID
		targetUsers = []string{req.ChatID}
		h.logger.Info("Sending to single chat",
			logger.String("chat_id", req.ChatID))
	} else {
		h.logger.Warn("No recipients specified in request")
		http.Error(w, "Either chat_id or user_ids must be specified", http.StatusBadRequest)
		return
	}

	// Send messages
	successCount, failCount := h.telegramClient.SendMessageToMultiple(req.BotToken, req.Message, targetUsers)

	// Prepare response
	response := models.SendMessageResponse{
		Success:      true,
		TotalSent:    len(targetUsers),
		SuccessCount: successCount,
		FailCount:    failCount,
	}

	h.logger.Info("Send message request completed",
		logger.Int("total", len(targetUsers)),
		logger.Int("success", successCount),
		logger.Int("failed", failCount))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Health handles GET requests for health check
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response := models.HealthResponse{
		Status:  "healthy",
		Service: "telegram-api-bridge",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
