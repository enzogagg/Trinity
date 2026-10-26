package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// pinger : ce dont l'API a besoin de la base (implémenté par *pgxpool.Pool, simulé dans les tests)
type pinger interface {
	Ping(ctx context.Context) error
}

func newMux(db pinger) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/health/db", dbHealthHandler(db))
	return mux
}

// healthHandler : liveness, ne dépend pas de la base (utilisé par le healthcheck Docker)
func healthHandler(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// dbHealthHandler : readiness, vérifie que PostgreSQL répond
func dbHealthHandler(db pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := db.Ping(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"database": "down"})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"database": "ok"})
	}
}

func writeJSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}
