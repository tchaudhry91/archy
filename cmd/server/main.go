package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/tchaudhry91/archy/service"

	"github.com/tchaudhry91/archy/service/store"

	"github.com/go-kit/kit/log"

	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("archy-svc", flag.ExitOnError)
	var (
		bindAddr      = fs.String("bind-addr", "127.0.0.1:15999", "Address to bind on")
		dbConn        = fs.String("db-conn", "", "MongoDB connection String")
		signingSecret = fs.String("signing-secret", "", "Secret to sign tokens with")
	)
	err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("ARCHY_SERVICE"))
	if err != nil {
		panic(err)
	}

	// DB
	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Store
	db, err := store.NewMongoStore(*dbConn)
	if err != nil {
		panic(err)
	}

	// Service
	svc := service.NewServer(db, logger, mux.NewRouter(), *bindAddr, *signingSecret)

	shutdown := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Log("msg", "Starting server..", "bindAddr", *bindAddr)
		err = svc.Start()
		shutdown <- err
	}()

	select {
	case signalKill := <-interrupt:
		logger.Log("msg", fmt.Sprintf("Stopping Server: %s", signalKill.String()))
	case err := <-shutdown:
		logger.Log("error", err)
	}

	err = svc.Shutdown()
	if err != nil {
		logger.Log("error", err)
	}

}
