-- name: CreateOrder :one
INSERT INTO orders (user_id, total_amount, status_staus)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE order_id = $1;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY order_date DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status_staus = $2, updated_at = now()
WHERE order_id = $1
RETURNING *;


-- name: DeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1;
