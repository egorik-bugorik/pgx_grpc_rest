package api

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"rest_api_postgres_clean/internal/inventory"
	"strings"
	"sync"
	"time"
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
	inventory  *inventory.Service
	server     *http.Server
	middleware func(http.Handler) http.Handler
}

type grpcServer struct {
	inventory *inventory.Service
	server    *grpc.Server
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
	for _ = range cap(errCh) {
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

	handler := NewHttpServer(s.inventory)

	if s.middleware != nil {
		handler = s.middleware(handler)
	}

	h := http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: time.Second * 5,
	}

	log.Printf("Server is listening on ::: %v", addr)

	if err := h.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil

}

func (s *grpcServer) Run(ctx context.Context, addr string) error {

	var ls net.ListenConfig
	lis, err := ls.Listen(ctx, "tcp", addr)
	if err != nil {
		return fmt.Errorf("error while listen ::: ", err)
	}
	s.server = grpc.NewServer()

	reflection.Register(s.server)

	grpcInvent := InventoryGRPC{
		inventory: s.inventory,
	}

	RegisterInventoryServer(s.server, grpcInvent)

	fmt.Printf("Server lsitening on %v", addr)
	if err := s.server.Serve(lis); err != nil {
		return fmt.Errorf("Fail to serve : %w", err)
	}
	return nil
}

// :::: STOP SERVERS

func (s *Server) Shutdown(ctx context.Context) {
	s.stopFn.Do(func() {
		s.httpServer.Shutdown(ctx)
		s.grpcServer.Shutdown(ctx)
	})
}

func (s *grpcServer) Shutdown(ctx context.Context) {

}
