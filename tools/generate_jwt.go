package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := flag.String("secret", "", "JWT secret key (required)")
	days := flag.Int("days", 30, "Token validity in days (default: 30)")
	flag.Parse()

	if *secret == "" {
		fmt.Println("Error: JWT secret is required")
		fmt.Println("Usage: go run generate_jwt.go -secret YOUR_SECRET [-days 30]")
		return
	}

	// Create token with expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24 * time.Duration(*days)).Unix(),
		"iat": time.Now().Unix(),
	})

	// Sign the token
	tokenString, err := token.SignedString([]byte(*secret))
	if err != nil {
		fmt.Printf("Error generating token: %v\n", err)
		return
	}

	fmt.Println("JWT Token Generated Successfully!")
	fmt.Println("===============================")
	fmt.Printf("Token: %s\n", tokenString)
	fmt.Printf("Expires in: %d days\n", *days)
	fmt.Println("===============================")
	fmt.Println("\nUse this token in your API requests:")
	fmt.Printf(`{"JWT":"%s","TELEGRAM_BOT_TOKEN":"...","USER_IDS":[...],...}`, tokenString)
	fmt.Println()
}
