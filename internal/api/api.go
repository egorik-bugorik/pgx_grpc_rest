package api

import (
	"context"
	"errors"
	"fmt"
	"rest_api_postgres_clean/internal/inventory"
	"strings"
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

	ctx, cancel := context.WithCancel(ctx)
	//:::: error channel for starting server
	errCh := make(chan error, 2)

	//::: start goroutine SERVER (http,grps)

	go func() {
		err := s.httpServer.Run(ctx, s.HTTPaddr)
		if err != nil {
			err = fmt.Errorf(" http Server finished with error :::%w ", err)
		}
		errCh <- err

	}()
	go func() {
		err := s.grpcServer.Run(ctx, s.GRPCaddr)
		if err != nil {
			err = fmt.Errorf(" grpc Server finished with error :::%w ", err)
		}
		errCh <- err

	}()

	var stringErr []string

	// ::: HANDLE errors
	for i := range cap(errCh) {
		if err := <-errCh; err != nil {
			stringErr = append(stringErr, err.Error())

			// SHUTDOWN server in case of healthy context
			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}

	}
	cancel()

	err := errors.New(strings.Join(stringErr, ", "))
	return err

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
