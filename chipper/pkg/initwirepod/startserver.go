package initwirepod

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	chipperpb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/api/go/jdocspb"
	"github.com/digital-dream-labs/api/go/tokenpb"
	"github.com/digital-dream-labs/hugh/log"
	"github.com/kercre123/chipper/pkg/logger"
	chipperserver "github.com/kercre123/chipper/pkg/servers/chipper"
	jdocsserver "github.com/kercre123/chipper/pkg/servers/jdocs"
	tokenserver "github.com/kercre123/chipper/pkg/servers/token"
	wp "github.com/kercre123/chipper/pkg/wirepod/preqs"
	"github.com/kercre123/chipper/pkg/wirepod/speechrequest"
	"github.com/soheilhy/cmux"

	//	grpclog "github.com/digital-dream-labs/hugh/grpc/interceptors/logger"

	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
)

var CmuxRunning bool = false
var ServerRunning bool = false
var GrpcServer *grpcserver.Server
var GrpcListener net.Listener
var m cmux.CMux
var p *wp.Server

func serveOk(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ok")
	return
}

func httpServe(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ok:80", serveOk)
	mux.HandleFunc("/ok", serveOk)
	s := &http.Server{
		Handler: mux,
	}
	return s.Serve(l)
}

func grpcServe(l net.Listener) error {
	var err error
	GrpcServer, err = grpcserver.New(
		grpcserver.WithViper(),
		grpcserver.WithReflectionService(),
		grpcserver.WithInsecureSkipVerify(),
	)
	if err != nil {
		log.Fatal(err)
	}

	s, _ := chipperserver.New(
		chipperserver.WithIntentProcessor(p),
		chipperserver.WithKnowledgeGraphProcessor(p),
		chipperserver.WithIntentGraphProcessor(p),
	)

	tokenServer := tokenserver.NewTokenServer()
	jdocsServer := jdocsserver.NewJdocsServer()
	//jdocsserver.IniToJson()

	chipperpb.RegisterChipperGrpcServer(GrpcServer.Transport(), s)
	jdocspb.RegisterJdocsServer(GrpcServer.Transport(), jdocsServer)
	tokenpb.RegisterTokenServer(GrpcServer.Transport(), tokenServer)

	return GrpcServer.Transport().Serve(l)
}

func StopGrpc() {
	ServerRunning = false
	logger.Println("gRPC server stopped successfully")
}

func StartServer(sttInitFunc func() error, sttHandlerFunc func(speechrequest.SpeechRequest) (string, error), voiceProcessorName string) {
	// begin wirepod stuff
	var err error
	p, err = wp.New(sttInitFunc, sttHandlerFunc, voiceProcessorName)
	//go wpweb.StartWebServer()
	wp.InitHoundify()
	if err != nil {
		logger.Println("Error starting server")
		logger.Println(err)
		ServerRunning = false
		return
	}
	if !CmuxRunning {
		cert, err := tls.X509KeyPair([]byte(os.Getenv("DDL_RPC_TLS_CERTIFICATE")), []byte(os.Getenv("DDL_RPC_TLS_KEY")))
		if err != nil {
			logger.Println("Error starting server")
			logger.Println(err)
			ServerRunning = false
			return
		}
		listener, err := tls.Listen("tcp", ":"+os.Getenv("DDL_RPC_PORT"), &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			logger.Println("Error starting server")
			logger.Println(err)
			ServerRunning = false
			return
		}
		m = cmux.New(listener)
		httpListener := m.Match(cmux.HTTP1Fast())
		go httpServe(httpListener)
		GrpcListener = m.Match(cmux.HTTP2())
		go grpcServe(GrpcListener)
		CmuxServe(m)
		ServerRunning = true
	}
	fmt.Println("wire-prod-pod started successfully!")
}

func CmuxServe(m cmux.CMux) {
	if !CmuxRunning {
		logger.Println("Serving cmux")
		CmuxRunning = true
		go m.Serve()
	} else {
		logger.Println("Cmux already running")
	}
}
