package web

import (
	"net"
	"net/http"

	"github.com/clevergo/clevergo"
	"github.com/clevergo/log"
)

type Server struct {
	*http.Server
	router      *clevergo.Router
	middlewares []func(http.Handler) http.Handler
	logger      log.Logger
}

func NewServer(router *clevergo.Router, logger log.Logger) *Server {
	return &Server{
		Server: &http.Server{},
		router: router,
		logger: logger,
	}
}

func (srv *Server) Use(middlewares ...func(http.Handler) http.Handler) {
	srv.middlewares = append(srv.middlewares, middlewares...)
}

func (srv *Server) prepare() {
	srv.Handler = srv.router
	for i := len(srv.middlewares) - 1; i >= 0; i-- {
		srv.Handler = srv.middlewares[i](srv.Handler)
	}
}

func (srv *Server) ListenAndServe() error {
	srv.prepare()
	srv.logger.Infof("Server started on http://%s", srv.Addr)
	return srv.Server.ListenAndServe()
}

func (srv *Server) ListenAndServeTLS(certFile, keyFile string) error {
	srv.prepare()
	srv.logger.Infof("Server started on http://%s", srv.Addr)
	return srv.Server.ListenAndServeTLS(certFile, keyFile)
}

func (srv *Server) ListenAndServeUnix() error {
	srv.logger.Infof("Server started on https://%s", srv.Addr)
	l, err := net.Listen("unix", srv.Addr)
	if err != nil {
		return err
	}
	return srv.Serve(l)
}

func (srv *Server) Serve(l net.Listener) error {
	srv.prepare()
	return srv.Server.Serve(l)
}
