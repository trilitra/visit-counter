package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file found, using system env")
	}

	r := chi.NewRouter()

	r.Use(ipMiddleware)

	store := NewMemoryStore()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"health":"ok"}`))
	})

	r.Get("/whoami", func(w http.ResponseWriter, r *http.Request) {
		ip, ok := GetIP(r.Context())
		if !ok {
			http.Error(w, "ip not found", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ip":"`+ip+`"}`))

	})

	r.Get("/visit", func(w http.ResponseWriter, r *http.Request) {
		ip, ok := GetIP(r.Context())
		if !ok {
			http.Error(w, "ip not found", http.StatusInternalServerError)
			return
		}

		visits, err := store.Inc(ip)
		if err != nil {
			http.Error(w, "error update ip visits", http.StatusInternalServerError)
			return
		}
		n := strconv.Itoa(visits)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ip":"` + ip + `", "visits":` + n + `}`))
	})

	r.Get("/stats", func(w http.ResponseWriter, r *http.Request) {
		count, err := store.UniqueCount()
		if err != nil {
			http.Error(w, "error upload stats", http.StatusInternalServerError)
			return
		}
		n := strconv.Itoa(count)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"unique_ips":` + n + `}`))
	})

	log.Fatal(http.ListenAndServe(":8080", r))

}
