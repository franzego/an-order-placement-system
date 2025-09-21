package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/franzego/notification-service/events"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Notification struct {
	rdb            *redis.Client
	apiKey         string
	processedCount int64
	errorCount     int64
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	//PoolSize int
}

func NewNotificationService(config RedisConfig) (*Notification, error) {
	// Load environment variables once during initialization
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	apiKey := os.Getenv("apikey")
	if apiKey == "" {
		return nil, fmt.Errorf("SendGrid API key not found in environment variables")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Notification{
		rdb:    rdb,
		apiKey: apiKey,
	}, nil
}

type OrderedEvent struct {
	OrderID     int                    `json:"orderid"`
	UserID      int                    `json:"userid"`
	TotalAmount string                 `json:"totalamount"`
	Status      string                 `json:"status"`
	Items       []events.OrderItemData `json:"items"`
}

func (n *Notification) HandleEvent(ctx context.Context, msg []byte) error {
	log.Printf("DEBUG: Received raw message: %s", string(msg))

	var event events.Event
	if err := json.Unmarshal(msg, &event); err != nil {
		n.errorCount++
		log.Printf("ERROR: Failed to unmarshal event: %v", err)
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	log.Printf("DEBUG: Event type: %s, ID: %s", event.Type, event.ID)
	log.Printf("DEBUG: Event data: %+v", event.Data)

	// Handle different event types
	switch event.Type {
	case events.OrderCreatedEventType:
		return n.handleOrderCreatedEvent(ctx, event)
	case events.OrderStatusUpdatedEventType:
		return n.handleOrderUpdatedEvent(ctx, event)
	default:
		n.errorCount++
		log.Printf("ERROR: Unknown event type: %s", event.Type)
		return fmt.Errorf("unknown event type: %s", event.Type)
	}
}

func (n *Notification) handleOrderCreatedEvent(ctx context.Context, event events.Event) error {
	// Extract the order data from the event
	log.Printf("DEBUG: Event type: %T", event)
	log.Printf("DEBUG: Event data: %+v", event.Data)
	log.Printf("DEBUG: Event data type: %T", event.Data)

	if event.Data == nil {
		n.errorCount++
		return fmt.Errorf("event data is nil")
	}

	orderDataBytes, err := json.Marshal(event.Data)
	if err != nil {
		n.errorCount++
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	var orderEvent events.OrderCreatedEvent
	if err := json.Unmarshal(orderDataBytes, &orderEvent); err != nil {
		n.errorCount++
		log.Printf("ERROR: Failed to unmarshal order created event: %v", err)
		log.Printf("ERROR: Raw order data bytes: %s", string(orderDataBytes))
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	log.Printf("DEBUG: Processing order created event - OrderID: %d, UserID: %d", orderEvent.OrderID, orderEvent.UserID)
	log.Printf("DEBUG: Order event details - TotalAmount: %s, Status: %s, Items count: %d",
		orderEvent.TotalAmount, orderEvent.Status, len(orderEvent.Items))

	// Extract the actual user ID from the UserID field (which is coming as a map)

	userEmailKey := fmt.Sprintf("user:%d", orderEvent.UserID)
	userEmail, err := n.rdb.HGet(ctx, userEmailKey, "email").Result()

	log.Printf("DEBUG: Looking for email with key: %s", userEmailKey)
	//userData, err := n.rdb.HGetAll(ctx, userKey).Result()
	//userEmail, err := n.rdb.HGetAll(ctx, userEmailKey).Result()
	if err != nil {
		n.errorCount++
		log.Printf("ERROR: Failed to get email for UserID %d from key %s: %v", orderEvent.UserID, userEmailKey, err)
		return fmt.Errorf("failed to get email for UserID %d: %w", orderEvent.UserID, err)
	}
	log.Printf("DEBUG: Found email for UserID %d: %s", orderEvent.UserID, userEmail)
	/*if len(userData) == 0 {
		n.errorCount++
		log.Print("DEBUG: Checking to see if the programs reaches here before exiting2")
		log.Printf("DEBUG: User data is empty for key: %s", userKey)
		return fmt.Errorf("user not found for userID %d", orderEvent.UserID)
	}*/

	// Step 2: Get user data using the email
	/*userKey := fmt.Sprintf("user:%s", userEmail)
	log.Printf("DEBUG: Looking for user data with key: %s", userKey)

	userData, err := n.rdb.HGetAll(ctx, userKey).Result()
	if err != nil {
		n.errorCount++
		log.Printf("ERROR: Redis error while getting user data from key %s: %v", userKey, err)
		return fmt.Errorf("redis error while getting user data: %w", err)
	}*/

	if len(userEmailKey) == 0 {
		n.errorCount++
		log.Printf("ERROR: User data is empty for key: %s", userEmailKey)
		return fmt.Errorf("user not found for email %s", userEmail)
	}

	log.Printf("DEBUG: Retrieved user data: %+v", userEmailKey)

	// Create email content
	plainTextContent := generatePlainTextEmail(orderEvent)
	htmlContent := generateHTMLEmail(orderEvent)

	// Prepare email
	from := mail.NewEmail("One and Only David Store", "davidenenama@gmail.com")
	subject := fmt.Sprintf("Order Confirmation #%d", orderEvent.OrderID)
	to := mail.NewEmail("", userEmail)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	// Send email
	client := sendgrid.NewSendClient(n.apiKey)
	response, err := client.Send(message)
	if err != nil {
		n.errorCount++
		log.Printf("ERROR: Failed to send email via SendGrid: %v", err)
		return fmt.Errorf("failed to send email via SendGrid: %w", err)
	}

	// Check response status
	if response.StatusCode >= 400 {
		n.errorCount++
		log.Printf("ERROR: SendGrid API error: status=%d, body=%s", response.StatusCode, response.Body)
		return fmt.Errorf("SendGrid API error: status=%d, body=%s", response.StatusCode, response.Body)
	}

	n.processedCount++
	log.Printf("SUCCESS: Email sent for order %d to user %d (status: %d) - Total processed: %d",
		orderEvent.OrderID, orderEvent.UserID, response.StatusCode, n.processedCount)
	return nil
}

func (n *Notification) handleOrderUpdatedEvent(ctx context.Context, event events.Event) error {
	// Extract the order data from the event
	orderDataBytes, err := json.Marshal(event.Data)
	if err != nil {
		n.errorCount++
		return fmt.Errorf("failed to marshal order data: %w", err)
	}

	var orderEvent events.UpdatedOrderEvent
	if err := json.Unmarshal(orderDataBytes, &orderEvent); err != nil {
		n.errorCount++
		return fmt.Errorf("failed to unmarshal order updated event: %w", err)
	}

	log.Printf("DEBUG: Processing order updated event - OrderID: %d, UserID: %d, Status: %s -> %s",
		orderEvent.OrderID, orderEvent.UserID, orderEvent.OldStatus, orderEvent.NewStatus)

	log.Printf("Order %d status updated from %s to %s", orderEvent.OrderID, orderEvent.OldStatus, orderEvent.NewStatus)

	n.processedCount++
	return nil
}

func generatePlainTextEmail(order events.OrderCreatedEvent) string {
	var itemsList []string
	for _, item := range order.Items {
		itemsList = append(itemsList, fmt.Sprintf("Product %d: %d x $%.2f", item.ProductID, item.Quantity, item.Price))
	}

	return fmt.Sprintf("Order #%d Confirmation\n\nStatus: %s\nTotal: %s\n\nItems:\n%s\n\nThank you for your order.",
		order.OrderID,
		order.Status,
		order.TotalAmount,
		strings.Join(itemsList, "\n"))
}

func generateHTMLEmail(order events.OrderCreatedEvent) string {
	var itemsList []string
	for _, item := range order.Items {
		itemsList = append(itemsList, fmt.Sprintf("Product %d: %d x $%.2f", item.ProductID, item.Quantity, item.Price))
	}

	itemsHTML := strings.Join(itemsList, "</li><li>")
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
		itemsHTML)
}

// GetMetrics returns basic processing metrics
func (n *Notification) GetMetrics() (processed, errors int64) {
	return n.processedCount, n.errorCount
}
