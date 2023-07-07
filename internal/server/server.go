package server

import "net/http"

type Server struct {
	httpServer *http.Server
}

func NewHTTPServer(addr string, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: handler,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}
