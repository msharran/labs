package server

func (s *Server) routes() {
	s.mux.HandleFunc("GET /hello/{name}", s.handleHelloGet())
	s.mux.HandleFunc("GET /admin", s.onlyAdmin(s.handleAdminGet()))
}
