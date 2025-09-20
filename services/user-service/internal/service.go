package internal

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/franzego/user-service/authentication"
	db "github.com/franzego/user-service/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Servicer interface {
	SignUp(ctx context.Context, username string, email string, passwordhash string, firstname pgtype.Text, lastname pgtype.Text) (string, error)
	Login(ctx context.Context, email string, passwordhash string) (string, error)
}

type Service struct {
	RepoServicer
	authentication.Auth
	rdb *redis.Client
}

func NewService(svc RepoServicer, redisClient *redis.Client) *Service {
	if svc == nil {
		return nil
	}
	return &Service{
		RepoServicer: svc,
		rdb:          redisClient,
	}
}
func (s *Service) SignUp(ctx context.Context, username string, email string,
	passwordhash string, firstname pgtype.Text, lastname pgtype.Text) (string, error) {
	hashpwd, err := bcrypt.GenerateFromPassword([]byte(passwordhash), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	user, err := s.RepoServicer.CreateUser(ctx, username, email, string(hashpwd), firstname, lastname)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate keys") {
			return "", errors.New("user already exists")
		}
		return "", err
	}
	// Cache user data in Redis
	userKey := fmt.Sprintf("user:%s", email)
	userData := map[string]interface{}{
		"id":         user.ID,
		"username":   user.Username,
		"email":      user.Email,
		"firstname":  user.FirstName.String,
		"lastname":   user.LastName.String,
		"is_active":  user.IsActive,
		"created_at": user.CreatedAt.Time.Format(time.RFC3339),
	}

	// Store in Redis with 24 hour expiration
	err = s.rdb.HMSet(ctx, userKey, userData).Err()
	if err != nil {
		log.Printf("Failed to cache user data: %v", err)
	} else {
		// Set expiration for the hash
		s.rdb.Expire(ctx, userKey, 24*time.Hour)
	}

	tok, err := s.Auth.CreateToken(user.Username)
	if err != nil {
		log.Fatal(err)
	}
	return tok, nil

}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	var user db.GetUserByEmailRow
	var err error

	// Try to get user data from Redis first
	userKey := fmt.Sprintf("user:%s", email)
	userData, err := s.rdb.HGetAll(ctx, userKey).Result()
	if err == nil && len(userData) > 0 {
		// Found in cache, but still need to get from DB for password verification
		user, err = s.RepoServicer.GetUserByEmail(ctx, email)
		if err != nil {
			return "", fmt.Errorf("invalid credentials")
		}
	} else {
		// Not found in cache or error, get from database
		user, err = s.RepoServicer.GetUserByEmail(ctx, email)
		if err != nil {
			return "", fmt.Errorf("invalid credentials")
		}

		// Cache the user data
		userData := map[string]interface{}{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"firstname":  user.FirstName.String,
			"lastname":   user.LastName.String,
			"is_active":  user.IsActive,
			"created_at": user.CreatedAt.Time.Format(time.RFC3339),
		}

		// Store in Redis with 24 hour expiration
		err = s.rdb.HMSet(ctx, userKey, userData).Err()
		if err != nil {
			log.Printf("Failed to cache user data: %v", err)
		} else {
			s.rdb.Expire(ctx, userKey, 24*time.Hour)
		}
	}

	// Compare password with stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	// Create authentication token
	tok, err := s.Auth.CreateToken(user.Username)
	if err != nil {
		return "", fmt.Errorf("error creating authentication token: %v", err)
	}
	return tok, nil
}
