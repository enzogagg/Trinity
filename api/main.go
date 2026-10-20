package main

import (
	"log"
	"net/http"
	"os"
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
	log.Println("api listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
