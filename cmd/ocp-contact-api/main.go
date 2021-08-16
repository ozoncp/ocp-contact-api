package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ozoncp/ocp-contact-api/internal/api"
	desc "github.com/ozoncp/ocp-contact-api/pkg/ocp-contact-api"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":8002"
)

func runGrpc(log zerolog.Logger) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen port %v: %v", grpcPort, err)
	}
	s := grpc.NewServer()
	desc.RegisterOcpContactApiServer(s, api.NewOcpContactApiServer(log))

	//fmt.Printf("Server listening on %s\n", *grpcEndpoint)
	if err := s.Serve(listen); err != nil {
		log.Fatal().Err(err).Msgf("failed to serve: %v", err)
	}

	return nil
}

func main() {
	fmt.Println("This is an Ozon Contact API")
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()

	if err := runGrpc(log); err != nil {
		log.Fatal().Err(err)
	}

}
