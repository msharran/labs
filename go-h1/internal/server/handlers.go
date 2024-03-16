package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleHelloGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			http.Error(w, "missing name", http.StatusBadRequest)
			return
		}
		w.Write([]byte("Hello, " + name + "!"))
	}
}

func (s *Server) handleAdminGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Admin!"))
	}
}

func (s *Server) handleSecretsGetKeys() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys := s.sm.Keys()
		err := json.NewEncoder(w).Encode(struct {
			Keys []string
		}{
			Keys: keys,
		})

		if err != nil {
			http.Error(w, "failed to encode keys", http.StatusInternalServerError)
			return
		}
	}
}

func (s *Server) handleSecretsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		secret, err := s.sm.Get(key)
		if err != nil {
			http.Error(w, "secret not found", http.StatusNotFound)
			return
		}

		w.Write([]byte(secret))
	}
}

func (s *Server) handleSecretsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		secret := r.Header.Get("X-Secret")

		if key == "" {
			http.Error(w, "missing key", http.StatusBadRequest)
			return
		}

		if secret == "" {
			http.Error(w, "missing secret", http.StatusBadRequest)
			return
		}

		err := s.sm.Set(key, secret)
		if err != nil {
			http.Error(w, "failed to set secret", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
