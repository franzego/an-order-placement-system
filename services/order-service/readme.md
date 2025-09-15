# 🛒 Order Service – E-Commerce Microservices  

This service is part of an **e-commerce microservices architecture**.  
It is responsible for:  

- Creating new orders and order items in the database.  
- Publishing **OrderCreated** events to Kafka so other services (store, notifications, billing) can react.  
- Exposing a REST API for order management.  

---

## 📂 Project Structure 
order-service/
├── cmd/ # entrypoint (main.go)
├── config/ # config for db, kafka, ports
├──
├── db/sqlc/ # auto-generated db access layer, migrations
├── internal/ # business logic (repo, service, handler,router)
├── events/events.go # defining the various events structs
├── kafka/ # defining the producer(publisher to publish the events)
├── service/ # defining the event service
├── tests/ # mock tests
├── sqlc.yaml # the configuration for sqlc
├── docker-compose.yml # to spin up kafka, zookeeper
└── README.md

---

## 🚀 Features  

- REST API endpoint to create orders.  
- PostgreSQL as the order database.  
- Kafka integration for event-driven architecture.  
- Clean architecture with clear separation:  
  - **repo** → DB queries, migrations 
  - **service** → business logic (DB, Kafka, order)  
  - **handler** → HTTP endpoints  

---

## ⚙️ Requirements  

- Go 1.21+  
- PostgreSQL 13+  
- Kafka (local or cluster)  
- Docker (optional for local dev)  

---

## 🔧 Setup  

1. **Clone the repo**  

```bash
git clone https://github.com/franzego/ecommerce-microservices/order-service.git
cd order-service

2. **Create a .env file** 
DB_URL=postgres://postgres:password@localhost:5432/nameofpostgresdb?sslmode=disable
KAFKA_BROKERS=localhost:9092
ORDER_TOPIC=orders
Port=8080

3. Run PostgreSQL + Kafka locally (via Docker Compose, for example).

4. **Start the service**
go run ./cmd/order-service


📡 API Endpoints
POST /orders

Create a new order.
Request body: 
{
  "user_id": 101,
  "args": [
    {
      "product_id": 1001,
      "quantity": 2,
      "price": 4.99
    },
    {
      "product_id": 1002,
      "quantity": 1,
      "price": 3.50
    }
  ]
}


📬 Kafka Event
When an order is created, an event is published to the orders topic:
{
	"id": "1757896460554384305",
	"type": "order.created",
	"timestamp": "2025-09-15T00:34:20.554387354Z",
	"version": "1.0",
	"data": {
		"orderid": 17,
		"userid": 0,
		"totalamount": "13.48",
		"status": "pending",
		"items": [
			{
				"product_id": 1001,
				"quantity": 2,
				"price": 4.99
			},
			{
				"product_id": 1002,
				"quantity": 1,
				"price": 3.5
			}
		]
	}
}
Other services (inventory, notification, billing) will consume this.

