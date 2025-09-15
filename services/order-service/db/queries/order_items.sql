-- name: AddOrderItem :one
INSERT INTO order_items (order_id, product_id, product_name, price, quantity)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListOrderItems :many
SELECT * FROM order_items
WHERE order_id = $1;

-- name: DeleteOrderItem :exec
DELETE FROM order_items
WHERE order_item_id = $1;
