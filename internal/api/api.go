package api

import (
	"context"
	"rest_api_postgres_clean/internal/inventory"
	"sync"
)

type Server struct {
	HTTPaddr string
	GRPCaddr string

	Service *inventory.Service

	httpServer *httpServer
	grpcServer *grpcServer

	stopFn sync.Once
}

type httpServer struct {
}

type grpcServer struct {
}

//:::: RUN SERVERS

func (s *Server) Run(ctx context.Context) error {

}
func (s *httpServer) Run(ctx context.Context, addr string) error {

}
func (s *grpcServer) Run(ctx context.Context, addr string) error {

}

// :::: STOP SERVERS

func (s *Server) Shutdown(ctx context.Context) {

}

func (s *httpServer) Shutdown(ctx context.Context) {

}

func (s *grpcServer) Shutdown(ctx context.Context) {

}
