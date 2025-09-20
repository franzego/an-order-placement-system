# 🛒 Store Microservices Project

This project is a microservices-based store application built with **Go**, **Postgres**, **Redis**, **Kafka**, and **Docker Compose**.  
It demonstrates separation of concerns across services and event-driven communication for notifications.

---

## 📦 Services Overview

### 🔹 Users Service (`localhost:8081`)
- Registers new users.  
- Stores user details in its own Postgres database (`usersdb`).  
- Caches user emails in Redis for quick lookup by Notifications Service.

**Endpoint:**
- `POST /signup`
  ```json
  { "email": "john@example.com", "password": "secret123" }
  ```
  Response:
  ```json
  { "user_id": 42, "email": "john@example.com" }
  ```

---

### 🔹 Orders Service (`localhost:8082`)
- Handles order creation.  
- Persists orders into its own Postgres database (`ordersdb`).  
- Publishes order events to Kafka for consumption.

**Endpoint:**
- `POST /order`
  ```json
  { "user_id": 42, "items": ["apple", "banana"], "total_amount": "12.99" }
  ```
  Response:
  ```json
  { "order_id": 123, "status": "pending" }
  ```

---

### 🔹 Store Service (`localhost:8083`)
- Acts as the storefront API (entry point for frontend).  
- Manages product catalog in `storedb`.  
- Proxies checkout calls to Orders Service.

**Endpoints:**
- `GET /products`
  ```json
  [
    { "id": 1, "name": "Apple", "price": "1.99" },
    { "id": 2, "name": "Banana", "price": "0.99" }
  ]
  ```

- `GET /products/{id}`
  ```json
  { "id": 1, "name": "Apple", "price": "1.99" }
  ```

- `POST /checkout`
  ```json
  { "user_id": 42, "items": [1, 2] }
  ```
  Response:
  ```json
  { "order_id": 123, "status": "pending" }
  ```

---

### 🔹 Notifications Service
- Listens to Kafka `orders` topic.  
- Looks up user email from Redis.  
- Sends order confirmation email via SendGrid.  
- No public API.

---

## 🗂 Databases

All services use a **single Postgres container** with multiple databases, initialized via `init.sql`:

```sql
CREATE DATABASE usersdb;
CREATE DATABASE ordersdb;
CREATE DATABASE storedb;
```

- Users Service → `usersdb`  
- Orders Service → `ordersdb`  
- Store Service → `storedb`  

---

## 🐳 Docker Compose Setup

### `docker-compose.yml`
```yaml
version: "3.9"

services:
  postgres:
    image: postgres:15
    container_name: postgres
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: pass
      POSTGRES_DB: usersdb
    ports:
      - "5433:5432"
    volumes:
      - ./init.sql:/docker-entrypoint-initdb.d/init.sql

  redis:
    image: redis:7
    container_name: redis
    ports:
      - "6379:6379"

  zookeeper:
    image: confluentinc/cp-zookeeper:7.4.0
    container_name: zookeeper
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
    ports:
      - "2181:2181"

  kafka:
    image: confluentinc/cp-kafka:7.4.0
    container_name: kafka
    depends_on:
      - zookeeper
    ports:
      - "9092:9092"
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092,PLAINTEXT_HOST://localhost:9092
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT
      KAFKA_INTER_BROKER_LISTENER_NAME: PLAINTEXT

  users-service:
    build: ./users-service
    container_name: users-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: user
      DB_PASSWORD: pass
      DB_NAME: usersdb
    depends_on:
      - postgres
      - redis
    ports:
      - "8081:8080"

  orders-service:
    build: ./orders-service
    container_name: orders-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: user
      DB_PASSWORD: pass
      DB_NAME: ordersdb
    depends_on:
      - postgres
      - kafka
    ports:
      - "8082:8080"

  store-service:
    build: ./store-service
    container_name: store-service
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: user
      DB_PASSWORD: pass
      DB_NAME: storedb
    depends_on:
      - postgres
    ports:
      - "8083:8080"

  notifications-service:
    build: ./notifications-service
    container_name: notifications-service
    environment:
      REDIS_HOST: redis
      REDIS_PORT: 6379
    depends_on:
      - redis
      - kafka
```

---

## 🚀 Running the Project

1. Ensure Docker and Docker Compose are installed.  
2. Start the stack:
   ```bash
   docker compose up --build
   ```
3. Access services:
   - Users Service → `http://localhost:8081`
   - Orders Service → `http://localhost:8082`
   - Store Service → `http://localhost:8083`
   - Postgres → `localhost:5433`
   - Redis → `localhost:6379`
   - Kafka → `localhost:9092`

---

## 🔗 Workflow Recap

1. **User signs up** → saved in `usersdb`, email cached in Redis.  
2. **Frontend hits Store Service** → fetches products / calls checkout.  
3. **Checkout** → Store Service proxies to Orders Service.  
4. **Orders Service** → saves to `ordersdb`, publishes to Kafka.  
5. **Notifications Service** → consumes Kafka event, looks up Redis email, sends confirmation.

---

## 📌 Notes
- Env vars (DB creds, Redis, Kafka, SendGrid key) should be placed in `.env` files per service.  
- This setup is for **local dev** only, not production.  

