package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/franzego/redisclient"
	"github.com/franzego/user-service/config"
	"github.com/franzego/user-service/internal"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.LoadConfig()
	dbconn2, err := pgxpool.New(context.Background(), cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer dbconn2.Close()

	rdb := redisclient.NewRedisService()
	if rdb == nil {
		log.Fatal("Failed to create Redis client")
	}

	repo := internal.NewRepoService(dbconn2)
	if repo == nil {
		log.Fatal("Failed to create repository service")
	}

	svc := internal.NewService(repo, rdb)
	if svc == nil {
		log.Fatal("Failed to create service")
	}

	handler := internal.NewHandler(svc)
	if handler == nil {
		log.Fatal("Failed to create handler service")
	}

	routes := internal.NewRouterService(handler)

	rou := mux.NewRouter()
	api := rou.PathPrefix("/api/v1").Subrouter()
	routes.RegisterRoutes(api)

	s := &http.Server{
		Handler:        rou,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		log.Printf("Starting server on %s", cfg.PORT)
		s.Addr = cfg.PORT
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = s.Shutdown(ctx)
	if err != nil {
		log.Printf("there was an errror %#v in shutting down", err)
	} else {
		log.Println("server shutdown gracefully")
	}
}
