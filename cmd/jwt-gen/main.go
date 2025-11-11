package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"telegram-bridge/config"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	// Define command-line flags
	expiresStr := flag.String("expires", "", "Token expiration duration (e.g., 24h, 7d, 30d, 0 for no expiration)")
	subject := flag.String("subject", "telegram-bridge-client", "JWT subject claim")
	issuer := flag.String("issuer", "telegram-bridge", "JWT issuer claim")
	showHelp := flag.Bool("help", false, "Show help message")

	flag.Parse()

	// Parse expiration duration with support for days
	var expiresIn time.Duration
	var err error
	if *expiresStr != "" && *expiresStr != "0" {
		// Check if it ends with 'd' for days
		if strings.HasSuffix(*expiresStr, "d") {
			daysStr := strings.TrimSuffix(*expiresStr, "d")
			days, parseErr := time.ParseDuration(daysStr + "h")
			if parseErr != nil {
				log.Fatalf("Invalid expiration duration: %v", parseErr)
			}
			expiresIn = days * 24
		} else {
			expiresIn, err = time.ParseDuration(*expiresStr)
			if err != nil {
				log.Fatalf("Invalid expiration duration: %v", err)
			}
		}
	}

	if *showHelp {
		fmt.Println("JWT Token Generator for Telegram Bridge")
		fmt.Println("\nUsage:")
		fmt.Println("  jwt-gen [options]")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  jwt-gen                           # Generate token with no expiration")
		fmt.Println("  jwt-gen -expires 24h              # Generate token that expires in 24 hours")
		fmt.Println("  jwt-gen -expires 168h             # Generate token that expires in 7 days")
		fmt.Println("  jwt-gen -subject \"my-app\"          # Generate token with custom subject")
		fmt.Println("\nThe JWT_SECRET is loaded from the .env file in the project root.")
		os.Exit(0)
	}

	// Load configuration to get JWT_SECRET
	cfg := config.Load()

	// Create JWT claims
	claims := jwt.MapClaims{
		"iss": *issuer,
		"sub": *subject,
		"iat": time.Now().Unix(),
	}

	// Add expiration if specified
	if expiresIn > 0 {
		claims["exp"] = time.Now().Add(expiresIn).Unix()
		fmt.Printf("Token will expire in: %s\n", expiresIn.String())
	} else {
		fmt.Println("Token will NOT expire (no expiration set)")
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		log.Fatalf("Failed to sign token: %v", err)
	}

	// Output the token
	separator := strings.Repeat("=", 80)
	fmt.Println("\n" + separator)
	fmt.Println("Generated JWT Token:")
	fmt.Println(separator)
	fmt.Println(tokenString)
	fmt.Println(separator)
	fmt.Println("\nToken Details:")
	fmt.Printf("  Subject: %s\n", *subject)
	fmt.Printf("  Issuer: %s\n", *issuer)
	fmt.Printf("  Issued At: %s\n", time.Now().Format(time.RFC3339))
	if expiresIn > 0 {
		fmt.Printf("  Expires At: %s\n", time.Now().Add(expiresIn).Format(time.RFC3339))
	} else {
		fmt.Printf("  Expires At: Never\n")
	}
	fmt.Println()
}
