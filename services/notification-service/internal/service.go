package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Notification struct {
	rdb *redis.Client
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	PoolSize int
}

func NewNotificationService(config RedisConfig) (*Notification, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: 10,
		MaxRetries:   3,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		DialTimeout:  15 * time.Second,
		PoolTimeout:  15 * time.Second,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Notification{
		rdb: rdb,
	}, nil
}

type OrderedEvent struct {
	OrderID     int      `json:"orderid"`
	UserID      int      `json:"userid"`
	TotalAmount string   `json:"totalamount"`
	Status      string   `json:"status"`
	Items       []string `json:"items"`
}

func (n *Notification) HandleEvent(ctx context.Context, msg []byte) error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	apiKey := os.Getenv("apikey")
	if apiKey == "" {
		return fmt.Errorf("SendGrid API key not found in environment variables")
	}

	var orderEvent OrderedEvent
	if err := json.Unmarshal(msg, &orderEvent); err != nil {
		return fmt.Errorf("failed to unmarshal order event: %w", err)
	}

	// Get user email from Redis with context timeout
	key := strconv.Itoa(orderEvent.UserID)
	userEmail, err := n.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return fmt.Errorf("user email not found for userID %d", orderEvent.UserID)
	} else if err != nil {
		return fmt.Errorf("redis error while getting user email: %w", err)
	}

	// Create email content
	plainTextContent := generatePlainTextEmail(orderEvent)
	htmlContent := generateHTMLEmail(orderEvent)

	// Prepare email
	from := mail.NewEmail("Your Store", "davidenenama@gmail.com")
	subject := fmt.Sprintf("Order Confirmation #%d", orderEvent.OrderID)
	to := mail.NewEmail("", userEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Send email
	client := sendgrid.NewSendClient(apiKey)
	response, err := client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	// Check response status
	if response.StatusCode >= 400 {
		return fmt.Errorf("SendGrid API error: status=%d, body=%s", response.StatusCode, response.Body)
	}

	log.Printf("Email sent successfully for order %d to user %d (status: %d)",
		orderEvent.OrderID, orderEvent.UserID, response.StatusCode)
	return nil

}

func generatePlainTextEmail(order OrderedEvent) string {
	return fmt.Sprintf("Order #%d Confirmation\n\nStatus: %s\nTotal: %s\n\nItems:\n%s\n\nThank you for your order.",
		order.OrderID,
		order.Status,
		order.TotalAmount,
		strings.Join(order.Items, "\n"))
}

func generateHTMLEmail(order OrderedEvent) string {
	itemsList := strings.Join(order.Items, "</li><li>")
	return fmt.Sprintf(`
<div style="font-family: sans-serif">
    <h2>Order #%d Confirmation</h2>
    <p>Status: %s<br>Total: %s</p>
    <p>Items:</p>
    <ul><li>%s</li></ul>
    <p>Thank you for your order.</p>
</div>`,
		order.OrderID,
		order.Status,
		order.TotalAmount,
		itemsList)
}
