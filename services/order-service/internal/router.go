package internal

import (
	"log"

	"github.com/gorilla/mux"
)

type routeService struct {
	*handle
}

func NewRouteService(h *handle) *routeService {
	return &routeService{handle: h}
}

func (rs *routeService) RegisterRoutes(r *mux.Router) {
	if rs.handle == nil {
		log.Fatal("handler is nil")
	}

	r.HandleFunc("/order", rs.handle.CreateOrder).Methods("POST")
	r.HandleFunc("/orders", rs.handle.ListOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", rs.handle.GetOrder).Methods("GET")
}
