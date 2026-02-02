package main

import (
	"context"
	"log"
	"net/http"

	"github.com/shahnajsc/OnePointLedger/backend/internal/api/auth"
	"github.com/shahnajsc/OnePointLedger/backend/internal/config"
	"github.com/shahnajsc/OnePointLedger/backend/internal/db"
	"github.com/shahnajsc/OnePointLedger/backend/internal/repo"
	"github.com/shahnajsc/OnePointLedger/backend/internal/service"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// test handler for sandbox callback URL
func opCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback received. Query params: " + query.Encode()))
}

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is missing")
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is missing")
	}

	ctx := context.Background()
	sqlDB, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	userRepo := repo.NewUserRepo(sqlDB)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/connect/op/callback", opCallbackHandler)

	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
