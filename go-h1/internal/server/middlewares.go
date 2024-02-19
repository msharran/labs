package server

import "net/http"

func (s *Server) onlyAdmin(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Admin") != "true" {
			http.Error(w, "not an admin", http.StatusForbidden)
			return
		}
		h(w, r)
	}
}
