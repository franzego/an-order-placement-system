-- name: CheckStock :one
SELECT available_quantity 
FROM products 
WHERE id = $1;

-- name: ReserveStock :exec
UPDATE products 
SET available_quantity = available_quantity - $1,
    reserved_quantity = reserved_quantity + $1
WHERE id = $2 AND available_quantity >= $1;

-- name: ReleaseStock :exec  
UPDATE products
SET available_quantity = available_quantity + $1,
    reserved_quantity = reserved_quantity - $1
WHERE id = $2 AND reserved_quantity >= $1;

-- name: GetProduct :one
SELECT id, product_name, available_quantity, price
FROM products
WHERE id = $1;