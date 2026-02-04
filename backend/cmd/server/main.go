package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/shahnajsc/OnePointLedger/backend/internal/api/auth"
	"github.com/shahnajsc/OnePointLedger/backend/internal/api/middleware"
	"github.com/shahnajsc/OnePointLedger/backend/internal/api/opconnect"
	"github.com/shahnajsc/OnePointLedger/backend/internal/api/user"
	"github.com/shahnajsc/OnePointLedger/backend/internal/config"
	"github.com/shahnajsc/OnePointLedger/backend/internal/db"
	"github.com/shahnajsc/OnePointLedger/backend/internal/opclient"
	"github.com/shahnajsc/OnePointLedger/backend/internal/repo"
	"github.com/shahnajsc/OnePointLedger/backend/internal/service"
	"github.com/joho/godotenv"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func opCallbackHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Callback received. Query params: " + query.Encode()))
}

func mustEnv(name, value string) {
	if value == "" {
		log.Fatalf("%s is missing", name)
	}
}

func main() {
	// Load config
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}
	cfg := config.Load()

	// Required basic env
	mustEnv("DATABASE_URL", cfg.DatabaseURL)
	mustEnv("JWT_SECRET", cfg.JWTSecret)

	// Required OP env (for /connect/op/start)
	mustEnv("OP_MTLS_BASE", cfg.OPMTLSBase)
	mustEnv("OP_AUTH_BASE", cfg.OPAuthBase)
	mustEnv("OP_CLIENT_ID", cfg.OPClientID)
	mustEnv("OP_CLIENT_SECRET", cfg.OPClientSecret)
	mustEnv("OP_API_KEY", cfg.OPAPIKey)
	mustEnv("OP_FAPI_FINANCIAL_ID", cfg.OPFAPIFinancialID)
	mustEnv("OP_REDIRECT_URI", cfg.OPRedirectURI)
	mustEnv("OP_QWAC_CERT_PATH", cfg.OPQWACCertPath)
	mustEnv("OP_QWAC_KEY_PATH", cfg.OPQWACKeyPath)
	mustEnv("OP_QSEAL_KEY_PATH", cfg.OPQSEALKeyPath)
	mustEnv("OP_QSEAL_KID", cfg.OPQSEALKid)

	log.Println("OP client_id:", cfg.OPClientID) // for Debug

	// Audience for OP request JWT (optional?)
	opAud := os.Getenv("OP_REQUEST_AUD")
	if opAud == "" {
		opAud = cfg.OPMTLSBase // reasonable default for sandbox
	}

	// DB connection (database/sql pool)
	ctx := context.Background()
	sqlDB, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Auth / User
	userRepo := repo.NewUserRepo(sqlDB)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := auth.NewHandler(authSvc)
	userHandler := user.NewHandler()

	// Middleware: to ensure protected route
	authMiddleware := middleware.JWTAuth(cfg.JWTSecret)

	// OP Connect dependencies: mTLS HTTP client using QWAC
	opHTTP, err := opclient.NewMTLSClient(cfg.OPQWACCertPath, cfg.OPQWACKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// OP Connect dependencies: AIS client (client credentials + /authorizations)
	ais := &opclient.AISClient{
		HTTP:            opHTTP,
		MTLSBase:        cfg.OPMTLSBase,
		ClientID:        cfg.OPClientID,
		ClientSecret:    cfg.OPClientSecret,
		APIKey:          cfg.OPAPIKey,
		FAPIFinancialID: cfg.OPFAPIFinancialID,
	}

	// OP Connect dependencies: Repo to store state -> authorizationId mapping
	opRepo := repo.NewOPConnectRepo(sqlDB)

	// OP Connect dependencies: Service creates auth intent + signs request JWT (QSEAL)
	opSvc, err := service.NewOPConnectService(
		ais,
		opRepo,
		cfg.OPAuthBase,
		cfg.OPRedirectURI,
		cfg.OPClientID,
		opAud,
		cfg.OPQSEALKeyPath,
		cfg.OPQSEALKid,
	)
	if err != nil {
		log.Fatal(err)
	}

	// OP Connect dependencies: HTTP handler
	opHandler := opconnect.NewHandler(opSvc)

	// Router
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/auth/signup", authHandler.Signup)
	mux.HandleFunc("/auth/login", authHandler.Login)
	mux.HandleFunc("/connect/op/callback", opCallbackHandler) // callback must be public because OP redirects without JWT

	// Routes: protected
	mux.Handle("/me", authMiddleware(http.HandlerFunc(userHandler.Me)))
	mux.Handle("/connect/op/start", authMiddleware(http.HandlerFunc(opHandler.Start)))

	// Backend
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
