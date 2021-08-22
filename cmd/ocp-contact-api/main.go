package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-contact-api/internal/repo"
	"net"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/ozoncp/ocp-contact-api/internal/api"
	"github.com/ozoncp/ocp-contact-api/internal/config"
	desc "github.com/ozoncp/ocp-contact-api/pkg/ocp-contact-api"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

func runGrpc(config *config.Config, log zerolog.Logger) error {
	listen, err := net.Listen("tcp", config.Grpc.Address)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to listen port %v: %v", config.Grpc.Address, err)
	}

	// Connect to the database
	dataSourceName := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		config.Database.SSL)

	db, err := sqlx.Connect(config.Database.Driver, dataSourceName)
	if err != nil {
		log.Error().Err(err).Msg("db connect failed")
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal().Err(err).Msg("db ping failed")
	}

	// start Grpc server
	s := grpc.NewServer()
	newRepo := repo.NewRepo(db)
	desc.RegisterOcpContactApiServer(s, api.NewOcpContactApiServer(newRepo, log))

	if err := s.Serve(listen); err != nil {
		log.Fatal().Err(err).Msgf("failed to serve: %v", err)
	}

	return nil
}

func main() {
	const configPath string = "./config.yml"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).
		With().Timestamp().Logger()

	cfg, err := config.Read(configPath)
	if err != nil {
		log.Fatal().Msgf("read configuration file by path \"%v\" failed", configPath)
		return
	}

	if err := runGrpc(cfg, log); err != nil {
		log.Fatal().Err(err)
	}

}
