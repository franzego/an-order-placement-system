package internal

import (
	"log"

	"github.com/gorilla/mux"
)

type RouteService struct {
	handle *Handler
}

func NewRouterService(hand *Handler) *RouteService {
	return &RouteService{handle: hand}
}

func (rs *RouteService) RegisterRoutes(r *mux.Router) {
	if rs.handle == nil {
		log.Fatal("handler is nil")
	}

	r.HandleFunc("/signup", rs.handle.Signup).Methods("POST")
	r.HandleFunc("/login", rs.handle.Login).Methods("POST")

}
