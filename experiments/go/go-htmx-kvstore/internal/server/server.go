package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

type Server struct {
	e    *echo.Echo
	db   *gorm.DB
	addr string
}

type ServerOpts struct {
	Port       string
	DBFileName string
	LocalHost  bool
}

func WithPort(port string) func(*ServerOpts) {
	return func(o *ServerOpts) {
		o.Port = port
	}
}

func WithDBFileName(dbFileName string) func(*ServerOpts) {
	return func(o *ServerOpts) {
		o.DBFileName = dbFileName
	}
}

func WithLocalHost() func(*ServerOpts) {
	return func(o *ServerOpts) {
		o.LocalHost = true
	}
}

func New(opts ...func(*ServerOpts)) (*Server, error) {
	o := &ServerOpts{ // default values
		Port:       "1323",
		DBFileName: "kvstore.sqlite3",
	}
	for _, f := range opts {
		f(o)
	}

	addr := ":" + o.Port
	if o.LocalHost {
		addr = "localhost" + addr
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())

	r, err := newRenderer()
	if err != nil {
		return nil, err
	}
	e.Renderer = r

	db, err := setupDB(o.DBFileName)
	if err != nil {
		return nil, err
	}

	s := &Server{e: e, db: db, addr: addr}
	s.setupRoutes()
	return s, nil
}

func (s *Server) ListenAndServe() error {
	return s.e.Start(s.addr)
}
