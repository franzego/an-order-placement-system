package internal

type OrderItemInput struct {
	ProductID   int64
	ProductName string
	Price       float64
	Quantity    int32
}

type OrderItemDTO struct {
	ID          int64
	ProductID   int64
	ProductName string
	Price       float64
	Quantity    int32
	Subtotal    float64
}
