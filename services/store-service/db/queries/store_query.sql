-- name: CheckStockAvailability :many
SELECT 
    id,
    product_name,
    available_quantity,
    reserved_quantity,
    (available_quantity >= $2) as is_sufficient
FROM store 
WHERE id = ANY($1::bigint[])
FOR UPDATE; -- Lock rows during transaction

-- name: ReserveStock :one
UPDATE store 
SET 
    available_quantity = available_quantity - $2,
    reserved_quantity = reserved_quantity + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
  AND available_quantity >= $2
RETURNING id, product_name, available_quantity, reserved_quantity;


-- name: ReleaseReservedStock :one
UPDATE store 
SET 
    available_quantity = available_quantity + $2,
    reserved_quantity = reserved_quantity - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
  AND reserved_quantity >= $2
RETURNING id, product_name, available_quantity, reserved_quantity;

-- name: ConfirmStockSale :one
UPDATE store 
SET 
    reserved_quantity = reserved_quantity - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 
  AND reserved_quantity >= $2
RETURNING id, product_name, available_quantity, reserved_quantity;

-- name: BulkCheckStock :many
SELECT 
    id,
    product_name,
    available_quantity,
    CASE 
        WHEN available_quantity >= $2 THEN true 
        ELSE false 
    END as can_fulfill
FROM store 
WHERE id = ANY($1::bigint[])
FOR UPDATE;


-- name: RestoreStock :one
UPDATE store 
SET 
    available_quantity = available_quantity + $2,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id, product_name, available_quantity, reserved_quantity;