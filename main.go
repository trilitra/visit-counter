package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type whoamiResponse struct {
	IP string `json:"ip"`
}

type visitResponse struct {
	IP     string `json:"ip"`
	Visits int64  `json:"visits"`
}

type statsResponse struct {
	UniqueIPs int64 `json:"unique_ips"`
}

type healthResponse struct {
	Status string `json:"status"`
}

type handler struct {
	store Store
}

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 10 * time.Second
	idleTimeout  = 60 * time.Second
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found, using system env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	r.Use(ipMiddleware)

	store := NewStore()

	h := &handler{store: store}

	r.Get("/health", handleHealth)

	r.Get("/whoami", handleWhoami)

	r.Get("/visit", h.handleVisits)

	r.Get("/stats", h.handleStats)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	err = srv.ListenAndServe()
	if err != nil {
		slog.Error("Error starting server:", "err", err)
		os.Exit(1)
	}
}

func NewStore() Store {
	storeType := os.Getenv("STORE_TYPE")
	if storeType == "" {
		storeType = "memory"
	}

	switch storeType {
	case "redis":
		addr := os.Getenv("REDIS_ADDR")
		if addr == "" {
			addr = "localhost:6379"
		}

		rdb := redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		})

		err := rdb.Ping(context.Background()).Err()
		if err != nil {
			slog.Error("Redis connection failed", "err", err)
		}

		return NewRedisStore(rdb)

	default:
		return NewMemoryStore()
	}
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(healthResponse{Status: "healthy"})
	if err != nil {
		slog.Error("Error encoding health response", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func handleWhoami(w http.ResponseWriter, r *http.Request) {
	ip, ok := GetIP(r.Context())
	if !ok {
		slog.Info("IP not found")
		http.Error(w, "ip not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(whoamiResponse{IP: ip})
	if err != nil {
		slog.Error("Error encoding ip response", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) handleVisits(w http.ResponseWriter, r *http.Request) {
	ip, ok := GetIP(r.Context())
	if !ok {
		slog.Info("IP not found")
		http.Error(w, "ip not found", http.StatusInternalServerError)
		return
	}

	visits, err := h.store.Inc(ip)
	if err != nil {
		slog.Error("Error update ip visits:", "err", err)
		http.Error(w, "error update ip visits", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(visitResponse{IP: ip, Visits: visits})
	if err != nil {
		slog.Error("Error encoding visits response:", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *handler) handleStats(w http.ResponseWriter, _ *http.Request) {
	count, err := h.store.UniqueCount()
	if err != nil {
		slog.Error("Error upload stats", "err", err)
		http.Error(w, "error upload stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(statsResponse{UniqueIPs: count})
	if err != nil {
		slog.Error("Error encoding stats response", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
