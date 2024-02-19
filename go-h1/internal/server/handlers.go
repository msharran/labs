package server

import "net/http"

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
