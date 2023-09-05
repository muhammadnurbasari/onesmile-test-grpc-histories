package main

import (
	"fmt"
	"net"
	"os"

	logkit "github.com/go-kit/kit/log"
	"github.com/joho/godotenv"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/db"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/endpointHistory"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/serviceHistory"
	"github.com/muhammadnurbasari/onesmile-test-grpc-histories/transport"
	"github.com/muhammadnurbasari/onesmile-test-protobuffer/proto/generate"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("config/.env")

	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")
	PORT := os.Getenv("PORT")

	log.Info().Msg(DB_HOST)
	log.Info().Msg(DB_PORT)
	log.Info().Msg(DB_USER)
	log.Info().Msg(DB_PASSWORD)
	log.Info().Msg(DB_NAME)
	log.Info().Msg(PORT)

	db, err := db.ConnPostgres(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)

	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(1)
	}

	fmt.Println("Succesfully connected")

	var logger logkit.Logger
	{
		logger = logkit.NewLogfmtLogger(os.Stderr)
		logger = logkit.With(logger, "ts", logkit.DefaultTimestampUTC)
		logger = logkit.With(logger, "caller", logkit.DefaultCaller)
	}

	var (
		service    = serviceHistory.NewServiceHistory(db)
		endpoints  = endpointHistory.NewEndpointHistory(service)
		grpcServer = transport.NewGrpcServer(endpoints, logger)
	)

	srv := grpc.NewServer()
	generate.RegisterTransactionsServer(srv, grpcServer)

	log.Info().Msg("Starting RPC server at " + ":" + PORT)

	l, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Error().Msg("could not listen to " + PORT + " error: " + err.Error())
		os.Exit(1)
	}

	if err := srv.Serve(l); err != nil {
		log.Error().Msg("could not listen to " + PORT + " error: " + err.Error())
		os.Exit(1)
	}
}
