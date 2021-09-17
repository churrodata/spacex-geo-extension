package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/churrodata/churro/pkg"
	"github.com/churrodata/spacex-geo-extension/internal/extension"
	pb "github.com/churrodata/spacex-geo-extension/rpc/extension"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	DEFAULT_PORT = ":10000"
)

func main() {
	zerolog.TimeFieldFormat = time.RFC822
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("spacex-geo-extension")

	debug := flag.Bool("debug", false, "debug flag")
	log.Info().Msg(fmt.Sprintf("debug set to %v", debug))
	serviceCertPath := flag.String("servicecert", "", "path to service cert files e.g. service.crt")
	dbCertPath := flag.String("dbcert", "", "path to database cert files (e.g. ca.crt)")

	flag.Parse()

	pipeline := os.Getenv("CHURRO_PIPELINE")
	if pipeline == "" {
		log.Error().Stack().Msg("CHURRO_PIPELINE env var is required")
		os.Exit(1)
	}
	ns := os.Getenv("CHURRO_NAMESPACE")
	if ns == "" {
		log.Error().Stack().Msg("CHURRO_NAMESPACE env var is required")
		os.Exit(1)
	}

	_, err := os.Stat(*dbCertPath)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in dbcertpath")
		os.Exit(1)
	}

	pi, err := pkg.GetPipeline()
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in getpipeline")
		os.Exit(1)
	}

	server := extension.NewExtensionServer(ns, true, *serviceCertPath, *dbCertPath, pi)
	creds, err := credentials.NewServerTLSFromFile(server.ServiceCreds.ServiceCrt, server.ServiceCreds.ServiceKey)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in getting server")
		os.Exit(1)
	}

	lis, err := net.Listen("tcp", DEFAULT_PORT)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in getting network")
		os.Exit(1)
	}

	s := grpc.NewServer(grpc.Creds(creds))

	pb.RegisterExtensionServer(s, server)
	if err := s.Serve(lis); err != nil {
		log.Error().Stack().Err(err).Msg("error in register")
		os.Exit(1)
	}

}
