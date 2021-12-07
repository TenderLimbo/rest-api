package models

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, router *gin.Engine) error {
	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
