package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// `api healthcheck` : utilisé par Docker (image distroless sans curl)
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		resp, err := http.Get("http://127.0.0.1:8080/api/health")
		if err != nil || resp.StatusCode != http.StatusOK {
			os.Exit(1)
		}
		return
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}
	// le pool se connecte à la demande : l'API démarre même si la base n'est pas encore prête
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("invalid DATABASE_URL: %v", err)
	}
	defer db.Close()

	log.Println("api listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", newMux(db)))
}
