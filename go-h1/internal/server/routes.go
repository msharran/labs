package server

func (s *Server) routes() {
	s.mux.HandleFunc("GET /hello/{name}", s.handleHelloGet())
	s.mux.HandleFunc("GET /admin", s.onlyAdmin(s.handleAdminGet()))

	s.mux.HandleFunc("GET /secrets/{key}", s.handleSecretsGet())
	s.mux.HandleFunc("POST /secrets/{key}", s.handleSecretsCreate())
	s.mux.HandleFunc("GET /secrets", s.handleSecretsGetKeys())
}
