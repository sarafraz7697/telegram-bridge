package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"telegram-bridge/config"
	"telegram-bridge/logger"
	"telegram-bridge/telegram"

	"github.com/golang-jwt/jwt/v5"
)

// BridgeRequest represents the incoming request structure
type BridgeRequest struct {
	JWT              string  `json:"JWT"`
	TelegramBotToken string  `json:"TELEGRAM_BOT_TOKEN"`
	UserIDs          []int64 `json:"USER_IDS"`
	Message          string  `json:"MESSAGE"`
	ImageBase64      string  `json:"IMAGE_BASE64,omitempty"`
	Type             string  `json:"TYPE"` // "message", "image", "image_and_message"
}

// Response represents the response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// TCPServer handles TCP connections and bridges to Telegram
type TCPServer struct {
	config *config.Config
	logger logger.Logger
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(cfg *config.Config, log logger.Logger) *TCPServer {
	return &TCPServer{
		config: cfg,
		logger: log,
	}
}

// Start starts the TCP server
func (s *TCPServer) Start() error {
	listener, err := net.Listen("tcp", ":"+s.config.Port)
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %w", err)
	}
	defer listener.Close()

	s.logger.Info("TCP server started", logger.String("port", s.config.Port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Error("Failed to accept connection", logger.String("error", err.Error()))
			continue
		}

		go s.handleConnection(conn)
	}
}

// handleConnection processes individual TCP connections
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	s.logger.Info("New connection", logger.String("remote_addr", remoteAddr))

	reader := bufio.NewReader(conn)

	// Read the incoming data
	data, err := reader.ReadString('\n')
	if err != nil {
		s.logger.Error("Failed to read from connection",
			logger.String("remote_addr", remoteAddr),
			logger.String("error", err.Error()))
		s.sendResponse(conn, false, "", "Failed to read data")
		return
	}

	data = strings.TrimSpace(data)
	s.logger.Debug("Received data", logger.String("remote_addr", remoteAddr), logger.Int("data_length", len(data)))

	// Parse the request
	var req BridgeRequest
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		s.logger.Error("Failed to parse request",
			logger.String("remote_addr", remoteAddr),
			logger.String("error", err.Error()))
		s.sendResponse(conn, false, "", "Invalid JSON format")
		return
	}

	// Validate JWT token
	if err := s.validateJWT(req.JWT); err != nil {
		s.logger.Warn("Unauthorized access attempt",
			logger.String("remote_addr", remoteAddr),
			logger.String("error", err.Error()))
		s.sendResponse(conn, false, "", "Unauthorized: Invalid JWT token")
		return
	}

	// Validate required fields
	if req.TelegramBotToken == "" {
		s.sendResponse(conn, false, "", "TELEGRAM_BOT_TOKEN is required")
		return
	}

	if len(req.UserIDs) == 0 {
		s.sendResponse(conn, false, "", "USER_IDS is required and cannot be empty")
		return
	}

	// Process based on type
	if err := s.processRequest(&req); err != nil {
		s.logger.Error("Failed to process request",
			logger.String("remote_addr", remoteAddr),
			logger.String("type", req.Type),
			logger.String("error", err.Error()))
		s.sendResponse(conn, false, "", err.Error())
		return
	}

	s.sendResponse(conn, true, "Message sent successfully", "")
}

// validateJWT validates the JWT token
func (s *TCPServer) validateJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWTSecret), nil
	})

	if err != nil {
		return fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// processRequest processes the bridge request and sends to Telegram
func (s *TCPServer) processRequest(req *BridgeRequest) error {
	// Create Telegram service dynamically with the bot token from the request
	tgService, err := telegram.NewService(req.TelegramBotToken, s.logger)
	if err != nil {
		return fmt.Errorf("failed to create Telegram service: %w", err)
	}

	switch req.Type {
	case "message":
		if req.Message == "" {
			return fmt.Errorf("MESSAGE is required for type 'message'")
		}
		return tgService.SendMessage(req.UserIDs, req.Message)

	case "image":
		if req.ImageBase64 == "" {
			return fmt.Errorf("IMAGE_BASE64 is required for type 'image'")
		}
		return tgService.SendImageBase64(req.UserIDs, req.ImageBase64, "")

	case "image_and_message":
		if req.ImageBase64 == "" || req.Message == "" {
			return fmt.Errorf("both IMAGE_BASE64 and MESSAGE are required for type 'image_and_message'")
		}
		return tgService.SendImageWithMessage(req.UserIDs, req.ImageBase64, req.Message)

	default:
		return fmt.Errorf("invalid TYPE: must be 'message', 'image', or 'image_and_message'")
	}
}

// sendResponse sends a JSON response back to the client
func (s *TCPServer) sendResponse(conn net.Conn, success bool, message string, errorMsg string) {
	resp := Response{
		Success: success,
		Message: message,
		Error:   errorMsg,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Failed to marshal response", logger.String("error", err.Error()))
		return
	}

	conn.Write(append(jsonResp, '\n'))
}
