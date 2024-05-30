package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"rest_api_postgres_clean/internal/api"
	"rest_api_postgres_clean/internal/database"
	"rest_api_postgres_clean/internal/inventory"
	"rest_api_postgres_clean/internal/postgres"
	"syscall"
	"time"
)

var (
	httpAddress = flag.String("http", "localhost:3000", "http addres of server is listen to... ")
	grpcAddress = flag.String("grps", "localhost:8082", "grps addres of server is listen to... ")
)

func main() {
	flag.Parse()

	// :::: setup LOGLEVEL

	logLvl, err := database.LogLevelFromEnv()
	if err != nil {
		return
	}
	// :::: setup POOL

	connStr := "host=localhost port=5432 user=postgres password=1 database=pgxtutorial"
	pool, err := database.NewPGXPool(context.Background(), connStr, &database.PGXStdLogger{}, logLvl)
	if err != nil {
		return
	}
	defer pool.Close()

	// :::: setup SERVER

	s := &api.Server{
		HTTPaddr: *httpAddress,
		GRPCaddr: *grpcAddress,
		Service:  inventory.NewService(postgres.NewDb(pool)),
	}

	//	 ::: CREATING ERROR CHANNEL

	chErr := make(chan error, 1)

	//STARTING SERVER IN GOROUTINE
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	go func() {
		chErr <- s.Run(ctx)

	}()

	//	SELECTING ON ctx Cancel or Server error

	select {
	//	CASE server error
	case <-chErr:

	//CASE context cancel\interrupt
	case <-ctx.Done():

		ctxShutdown, cancel := context.WithTimeout(context.Background(), time.Second*3)

		defer cancel()

		s.Shutdown(ctxShutdown)

		stop()

		err = <-chErr
	}

	if err != nil {
		log.Fatalf("Error while shut down!!! :::%w", err)

	}

}
