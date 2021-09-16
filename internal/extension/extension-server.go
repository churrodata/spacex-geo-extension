package extension

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/churrodata/churro/api/extract"
	extractapi "github.com/churrodata/churro/api/extract"
	"github.com/churrodata/churro/api/v1alpha1"
	"github.com/churrodata/churro/pkg/config"
	pb "github.com/churrodata/spacex-geo-extension/rpc/extension"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var DBCertPath string
var ServiceCertPath string

const (
	DEFAULT_PORT = ":10000"
)

type Server struct {
	Pi           v1alpha1.Pipeline
	ServiceCreds config.ServiceCredentials
	DBCreds      config.DBCredentials
	UserDBCreds  config.DBCredentials
}

// NewExtensionServer constructs a extension server based on the passed
// configuration, a pointer to the server is returned
func NewExtensionServer(ns string, debug bool, serviceCertPath string, dbCertPath string, pipeline v1alpha1.Pipeline) *Server {
	s := Server{
		Pi: pipeline,
	}

	err := s.SetupCredentials(ns, serviceCertPath, dbCertPath)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error getting credentials")
		os.Exit(1)
	}

	return &s
}

func (s *Server) Ping(ctx context.Context, request *pb.PingRequest) (response *pb.PingResponse, err error) {
	if false {
		return nil, status.Errorf(codes.InvalidArgument,
			"something is not right")
	}

	log.Info().Msg("ping called")

	return &pb.PingResponse{}, nil
}

func UnimplementedExtensionServer() {
}

func (s *Server) SetupCredentials(ns, serviceCertPath, dbCertPath string) (err error) {
	s.ServiceCreds = config.ServiceCredentials{
		ServiceCrt: serviceCertPath + "/service.crt",
		ServiceKey: serviceCertPath + "/service.key",
	}

	// churro-ctl uses the database root credentials
	s.DBCreds = config.DBCredentials{
		CACertPath:      dbCertPath + "/ca.crt",
		CAKeyPath:       dbCertPath + "/ca.key",
		SSLRootKeyPath:  dbCertPath + "/client.root.key",
		SSLRootCertPath: dbCertPath + "/client.root.crt",
		SSLKeyPath:      dbCertPath + "/client." + ns + ".key",
		SSLCertPath:     dbCertPath + "/client." + ns + ".crt",
	}

	s.UserDBCreds = config.DBCredentials{
		CACertPath:      dbCertPath + "/ca.crt",
		CAKeyPath:       dbCertPath + "/ca.key",
		SSLRootKeyPath:  dbCertPath + "/client.root.key",
		SSLRootCertPath: dbCertPath + "/client.root.crt",
		SSLKeyPath:      dbCertPath + "/client." + ns + ".key",
		SSLCertPath:     dbCertPath + "/client." + ns + ".crt",
	}

	err = s.ServiceCreds.Validate()
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in svccreds validate")
		return err
	}

	err = s.DBCreds.Validate()
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in dbcreds validate")
		return err
	}

	return nil
}

func (s *Server) Push(ctx context.Context, request *pb.PushRequest) (response *pb.PushResponse, err error) {

	loaderMsg := extractapi.LoaderMessage{}
	loaderMsg.Key = request.Key
	loaderMsg.Metadata = request.Metadata
	loaderMsg.DataFormat = request.DataFormat
	//	b := request.LoaderMessageString
	log.Info().Msg(fmt.Sprintf("got key %d dataformat %s len Metadata %d from Push function", loaderMsg.Key, loaderMsg.DataFormat, len(loaderMsg.Metadata)))

	var jsonStruct extractapi.RawFormat
	err = json.Unmarshal(loaderMsg.Metadata, &jsonStruct)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error in unmarshal")
		return nil, err
	}

	recCount := len(jsonStruct.Records)

	if recCount > 0 {
		err := s.processRecords(jsonStruct.Records)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error in unmarshal")
			return nil, err
		}

	}

	return &pb.PushResponse{}, nil
}

func (s Server) processRecords(records []extract.GenericRow) (err error) {
	url := s.DBCreds.GetDBConnectString(s.Pi.Spec.DataSource)
	log.Info().Msg("url " + url)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Error().Stack().Err(err).Msg("error connecting to the database")
		return err
	}
	for i, r := range records {
		log.Info().Msg(fmt.Sprintf("record %d key %d\n", i, r.Key))
		err := s.updateReverseGeo(db, r.Key)
		if err != nil {
			log.Error().Stack().Err(err).Msg("error updating geo code")
		}

	}
	return nil
}
