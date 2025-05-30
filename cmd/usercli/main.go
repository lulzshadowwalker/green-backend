package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/lulzshadowwalker/green-backend/internal/psql"
	"github.com/lulzshadowwalker/green-backend/internal/psql/db"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	fmt.Println("User Creation CLI")
	fmt.Println("-----------------")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Failed to read username: %v\n", err)
		os.Exit(1)
	}
	username = strings.TrimSpace(username)
	if username == "" {
		fmt.Println("Username cannot be empty.")
		os.Exit(1)
	}

	fmt.Print("Enter password: ")
	password, err := readPassword(reader)
	if err != nil {
		fmt.Printf("Failed to read password: %v\n", err)
		os.Exit(1)
	}
	password = strings.TrimSpace(password)
	if password == "" {
		fmt.Println("Password cannot be empty.")
		os.Exit(1)
	}

	pool, err := psql.Connect(psql.ConnectionParams{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	q := db.New(pool)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Failed to hash password: %v\n", err)
		os.Exit(1)
	}

	// Try to get user
	user, err := q.GetUserByUsername(ctx, username)
	if err != nil {
		// Insert new user
		_, err = pool.Exec(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", username, string(hash))
		if err != nil {
			fmt.Printf("Failed to insert user: %v\n", err)
			os.Exit(1)
		}
		user, err = q.GetUserByUsername(ctx, username)
		if err != nil {
			fmt.Printf("Failed to fetch user after insert: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Update password
		_, err = pool.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE username = $2", string(hash), username)
		if err != nil {
			fmt.Printf("Failed to update user password: %v\n", err)
			os.Exit(1)
		}
	}

	token, err := internal.GenerateJWT(int(user.ID), user.Username)
	if err != nil {
		fmt.Printf("Failed to generate JWT token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nUser created/updated successfully!")
	fmt.Println("JWT access token (save this somewhere safe):")
	fmt.Printf("\n%s\n\n", token)
	fmt.Println("Use this token as a Bearer token for API access.")
}

func readPassword(reader *bufio.Reader) (string, error) {
	pass, err := reader.ReadString('\n')
	return strings.TrimSpace(pass), err
}
