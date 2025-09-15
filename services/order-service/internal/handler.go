package internal

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	db "github.com/franzego/ecommerce-microservices/order-service/db/sqlc"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type handle struct {
	service Servicer
}

func NewHandleService(svc Servicer) *handle {
	if svc == nil {
		return nil
	}
	return &handle{service: svc}
}

// validation pkg from github to help me handle the validation
func (h *handle) Validate() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(h)
}

// handler to handle creating the order
func (h *handle) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	var req struct {
		Userid int                     `json:"userid" validate:"gt0"`
		Args   []db.AddOrderItemParams `json:"args" validate:"required,min=1"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	// for decoding purpose
	if len(req.Args) == 0 {
		log.Printf("ERROR: No items in request!")
		http.Error(w, "No items provided in request", http.StatusBadRequest)
		return
	}

	order, orderitem, err := h.service.CreateOrderWithItems(r.Context(), req.Userid, req.Args)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	type OrderResponse struct {
		Order     db.Order       `json:"order"`
		OrderItem []db.OrderItem `json:"orderitem"`
	}

	resp := &OrderResponse{
		Order:     order,
		OrderItem: orderitem,
	}

	w.Header().Set("Content-type", "application/json")

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order has been successfully created"))

}

//handler to List orders

func (h *handle) ListOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		http.Error(w, "Method not allowed", http.StatusBadRequest)
		return
	}
	var req struct {
		List db.ListOrdersParams
	}
	list, err := h.service.ListOrders(r.Context(), req.List)
	if err != nil {
		log.Print(err)
		http.Error(w, "error getting list of orders", http.StatusInternalServerError)
		return
	}
	type ListResponse struct {
		List []db.Order
	}
	resp := &ListResponse{
		List: list,
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of all orders"))

}

//func to get orders

func (h *handle) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	//to check for empty parameter
	ordidstr := vars["orderid"]
	if ordidstr == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	//to deal with the order id that will be passed in the uri
	id, err := strconv.ParseInt(vars["orderid"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID: Must be a number",
			http.StatusBadRequest)
		return
	}

	//var req struct {
	//	Orderid int32 `json:"orderid" validate:"gt0"`
	//}
	/*err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON"+err.Error(), http.StatusBadRequest)
		return
	}*/

	gottenorder, err := h.service.GetOrder(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	err = json.NewEncoder(w).Encode(gottenorder)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Get order with order id: " + ordidstr))

}
