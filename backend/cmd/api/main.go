package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"jobtrack/internal/database"
	"jobtrack/internal/httpapi"
	"jobtrack/internal/jobs"
)

func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	repo := jobs.NewRepository(db)
	jobHandler := jobs.NewHandler(repo)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", "GET")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("/api/jobs", jobHandler)
	mux.Handle("/api/jobs/", jobHandler)

	corsOrigins := os.Getenv("JOBTRACK_CORS_ORIGINS")
	if corsOrigins == "" {
		corsOrigins = "http://localhost:5173,http://127.0.0.1:5173"
	}

	allowedOrigins := strings.Split(corsOrigins, ",")
	handler := httpapi.CORSMiddleware(mux, allowedOrigins)

	addr := os.Getenv("JOBTRACK_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("JobTrack API listening on http://localhost%s", addr)
	log.Fatal(server.ListenAndServe())
}
