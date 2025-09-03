package server

func (s *Server) setupRoutes() {
	s.e.GET("/signup", s.handleViewSignup)
	s.e.POST("/signup", s.handleSignup)
	s.e.GET("/login", s.handleViewLogin)
	s.e.POST("/login", s.handleLogin)
	s.e.DELETE("/logout", s.handleLogout)

	s.e.GET("/kv", s.handleListKV)
	s.e.GET("/kv/new", s.handleNewKV)
	s.e.GET("/kv/:key/edit", s.handleEditKV)
	s.e.GET("/kv/:key/view", s.handleViewKV)
	s.e.POST("/kv", s.handleCreateKV)
	s.e.PUT("/kv/:key", s.handleUpdateKV)
	s.e.DELETE("/kv/:key", s.handleDeleteKV)
}
