package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	pb "github.com/digital-dream-labs/api/go/chipperpb"
	"github.com/digital-dream-labs/chipper/pkg/logger"
	"github.com/digital-dream-labs/chipper/pkg/server"
	"github.com/digital-dream-labs/chipper/pkg/voice_processors/wirepod-coqui"

	//	grpclog "github.com/digital-dream-labs/hugh/grpc/interceptors/log"

	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
	"github.com/digital-dream-labs/hugh/log"
)

var srv *grpcserver.Server
var grpcIsRunning bool = false

type chipperConfigStruct struct {
	Port           string `json:"port"`
	Cert           string `json:"cert"`
	Key            string `json:"key"`
	WeatherEnable  string `json:"weatherEnable"`
	WeatherKey     string `json:"weatherKey"`
	WeatherUnit    string `json:"weatherUnit"`
	HoundifyEnable string `json:"houndifyEnable"`
	HoundifyKey    string `json:"houndifyKey"`
	HoundifyID     string `json:"houndifyID"`
	STTService     string `json:"sttService"`
	PicovoiceKey   string `json:"picovoiceKey"`
}

func chipperAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case r.URL.Path == "/chipper/start_chipper":
		//	name := r.FormValue("name")
		if grpcIsRunning {
			fmt.Fprintf(w, "chipper already running")
		} else {
			go startServer()
			fmt.Fprintf(w, "chipper started")
		}
		return
	case r.URL.Path == "/chipper/stop_chipper":
		if !grpcIsRunning {
			fmt.Fprintf(w, "chipper already stopped")
		} else {
			stopServer()
			fmt.Fprintf(w, "chipper stopped")
		}
		return
	case r.URL.Path == "/chipper/restart_chipper":
		if !grpcIsRunning {
			go startServer()
			fmt.Fprintf(w, "chipper restarted")
		} else {
			stopServer()
			go startServer()
			fmt.Fprintf(w, "chipper restarted")
		}
		return
	case r.URL.Path == "/chipper/make_config":
		port := r.FormValue("port")
		cert := r.FormValue("cert")
		key := r.FormValue("key")
		weatherEnable := r.FormValue("weatherEnable")
		weatherKey := r.FormValue("weatherKey")
		weatherUnit := r.FormValue("weatherUnit")
		houndifyEnable := r.FormValue("houndifyEnable")
		houndifyKey := r.FormValue("houndifyKey")
		houndifyID := r.FormValue("houndifyID")
		if port == "" && cert == "" && key == "" {
			fmt.Fprintf(w, "port, cert, and key are required form values")
		}
		var chipperConfigReq chipperConfigStruct
		chipperConfigReq.Port = port
		chipperConfigReq.Cert = cert
		chipperConfigReq.Key = key
		chipperConfigReq.WeatherEnable = weatherEnable
		chipperConfigReq.WeatherKey = weatherKey
		chipperConfigReq.WeatherUnit = weatherUnit
		chipperConfigReq.HoundifyEnable = houndifyEnable
		chipperConfigReq.HoundifyKey = houndifyKey
		chipperConfigReq.HoundifyID = houndifyID
		chipperConfigBytes, _ := json.Marshal(chipperConfigReq)
		os.WriteFile("./chipperConfig.json", chipperConfigBytes, 0644)
		fmt.Fprintf(w, "chipper config created")
		return
	case r.URL.Path == "/chipper/edit_config":
		port := r.FormValue("port")
		cert := r.FormValue("cert")
		key := r.FormValue("key")
		weatherEnable := r.FormValue("weatherEnable")
		weatherKey := r.FormValue("weatherKey")
		weatherUnit := r.FormValue("weatherUnit")
		houndifyEnable := r.FormValue("houndifyEnable")
		houndifyKey := r.FormValue("houndifyKey")
		houndifyID := r.FormValue("houndifyID")
		var chipperConfigReq chipperConfigStruct
		chipperConfigOrig, err := os.ReadFile("./chipperConfig.json")
		if err != nil {
			fmt.Fprint(w, err.Error())
		}
		json.Unmarshal(chipperConfigOrig, &chipperConfigReq)
		if port != "" {
			chipperConfigReq.Port = port
		}
		if cert != "" {
			chipperConfigReq.Cert = cert
		}
		if key != "" {
			chipperConfigReq.Key = key
		}
		if weatherEnable != "" {
			chipperConfigReq.WeatherEnable = weatherEnable
		}
		if weatherKey != "" {
			chipperConfigReq.WeatherKey = weatherKey
		}
		if weatherUnit != "" {
			chipperConfigReq.WeatherUnit = weatherUnit
		}
		if houndifyEnable != "" {
			chipperConfigReq.HoundifyEnable = houndifyEnable
		}
		if houndifyKey != "" {
			chipperConfigReq.HoundifyKey = houndifyKey
		}
		if houndifyID != "" {
			chipperConfigReq.HoundifyID = houndifyID
		}
		chipperConfigBytes, _ := json.Marshal(chipperConfigReq)
		os.WriteFile("./chipperConfig.json", chipperConfigBytes, 0644)
		fmt.Fprintf(w, "chipper config edited")
		return
	case r.URL.Path == "/chipper/upload_ssh_key":
		r.ParseMultipartForm(32 << 20)
		var buf bytes.Buffer
		file, _, err := r.FormFile("file")
		if err != nil {
			fmt.Fprintf(w, "error")
			return
		}
		io.Copy(&buf, file)
		os.WriteFile("/tmp/sshKey", buf.Bytes(), 0600)
		fmt.Fprintf(w, "uploaded")
		return
	case r.URL.Path == "/chipper/setup_bot":
		botIP := r.FormValue("botIP")
		if _, err := os.Stat("/tmp/sshKey"); err != nil {
			fmt.Fprintf(w, "upload an ssh key first")
			return
		}
		cmdOutput, _ := exec.Command("/bin/bash", "./setupBot.sh", botIP, "/tmp/sshKey").Output()
		fmt.Fprintf(w, "Output: "+string(cmdOutput))
		return
	case r.URL.Path == "/chipper/get_logs":
		fmt.Fprintf(w, logger.LogList)
		return
	}
}

func main() {
	log.SetJSONFormat("2006-01-02 15:04:05")
	if os.Getenv("DDL_RPC_PORT") != "" {
		var chipperConfigReq chipperConfigStruct
		var shouldStartChipper bool = true
		origSourceBytes, _ := os.ReadFile("./source.sh")
		origSource := string(origSourceBytes)
		if strings.Contains(origSource, "../certs/cert.crt") {
			chipperConfigReq.Cert = "../certs/cert.crt"
			chipperConfigReq.Key = "../certs/cert.key"
		} else if strings.Contains(origSource, "./epod/ep.crt") {
			chipperConfigReq.Cert = "./epod/ep.crt"
			chipperConfigReq.Key = "./epod/ep.key"
		} else {
			logger.Logger("You must use the webserver to define where your certs are.")
			shouldStartChipper = false
		}
		chipperConfigReq.Port = os.Getenv("DDL_RPC_PORT")
		chipperConfigReq.WeatherEnable = os.Getenv("WEATHERAPI_ENABLED")
		chipperConfigReq.WeatherKey = os.Getenv("WEATHERAPI_KEY")
		chipperConfigReq.WeatherUnit = os.Getenv("WEATHERAPI_UNIT")
		chipperConfigReq.HoundifyEnable = os.Getenv("HOUNDIFY_ENABLED")
		chipperConfigReq.HoundifyKey = os.Getenv("HOUNDIFY_CLIENT_KEY")
		chipperConfigReq.HoundifyID = os.Getenv("HOUNDIFY_CLIENT_ID")
		chipperConfigBytes, _ := json.Marshal(chipperConfigReq)
		os.WriteFile("./chipperConfig.json", chipperConfigBytes, 0644)
		os.Rename("./source.sh", "old-source.sh")
		if shouldStartChipper {
			go startServer()
		}
	} else if _, err := os.Stat("./chipperConfig.json"); err == nil {
		go startServer()
	} else {
		logger.Logger("Use the webserver to setup and start chipper.")
	}
	var webPort string
	http.HandleFunc("/api/", server.ApiHandler)
	http.HandleFunc("/chipper/", chipperAPIHandler)
	fileServer := http.FileServer(http.Dir("./webroot"))
	http.Handle("/", fileServer)
	if os.Getenv("WEBSERVER_PORT") != "" {
		if _, err := strconv.Atoi(os.Getenv("WEBSERVER_PORT")); err == nil {
			webPort = os.Getenv("WEBSERVER_PORT")
		} else {
			logger.Logger("WEBSERVER_PORT contains letters, using default of 8080")
			webPort = "8080"
		}
	} else {
		webPort = "8080"
	}
	fmt.Printf("Starting webserver at port " + webPort + " (http://localhost:" + webPort + ")\n")
	if err := http.ListenAndServe(":"+webPort, nil); err != nil {
		log.Fatal(err)
	}
}

func stopServer() {
	srv.Stop()
	grpcIsRunning = false
}

func startServer() {
	if _, err := os.Stat("./chipperConfig.json"); err == nil {
		chipperConfigBytes, _ := os.ReadFile("./chipperConfig.json")
		var chipperConfig chipperConfigStruct
		json.Unmarshal(chipperConfigBytes, &chipperConfig)
		certBytes, _ := os.ReadFile(chipperConfig.Cert)
		certString := string(certBytes)
		keyBytes, _ := os.ReadFile(chipperConfig.Key)
		keyString := string(keyBytes)
		os.Setenv("DDL_RPC_PORT", chipperConfig.Port)
		os.Setenv("DDL_RPC_TLS_CERTIFICATE", certString)
		os.Setenv("DDL_RPC_TLS_KEY", keyString)
		os.Setenv("DDL_RPC_CLIENT_AUTHENTICATION", "NoClientCert")
		os.Setenv("WEATHERAPI_ENABLED", chipperConfig.WeatherEnable)
		os.Setenv("WEATHERAPI_KEY", chipperConfig.WeatherKey)
		os.Setenv("WEATHERAPI_UNIT", chipperConfig.WeatherUnit)
		os.Setenv("HOUNDIFY_ENABLED", chipperConfig.HoundifyEnable)
		os.Setenv("HOUNDIFY_CLIENT_KEY", chipperConfig.HoundifyKey)
		os.Setenv("HOUNDIFY_CLIENT_ID", chipperConfig.HoundifyID)
		os.Setenv("DEBUG_LOGGING", "true")
	}
	var err error
	srv, err = grpcserver.New(
		grpcserver.WithViper(),
		grpcserver.WithLogger(log.Base()),
		grpcserver.WithReflectionService(),

		grpcserver.WithUnaryServerInterceptors(
		//			grpclog.UnaryServerInterceptor(),
		),

		grpcserver.WithStreamServerInterceptors(
		//			grpclog.StreamServerInterceptor(),
		),
	)
	if err != nil {
		logger.Logger("Something is broken in the voice server config.")
		logger.Logger("GRPC server error: " + err.Error())
		logger.Logger("This can be solved via the webserver.")
		return
	}
	p, err := wirepod.New()
	var canGoOn bool = true
	if err != nil {
		logger.Logger("Something is broken in the voice server config.")
		logger.Logger("New wire-pod instance error: " + err.Error())
		logger.Logger("This can be solved via the webserver.")
		canGoOn = false
	}

	if canGoOn {
		s, _ := server.New(
			//server.WithLogger(log.Base()),
			server.WithIntentProcessor(p),
			server.WithKnowledgeGraphProcessor(p),
			server.WithIntentGraphProcessor(p),
		)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				logger.Logger("Interrupt detected, exiting")
				os.Exit(0)
			}
		}()

		pb.RegisterChipperGrpcServer(srv.Transport(), s)

		srv.Start()
		logger.Logger("Server started successfully!")
		grpcIsRunning = true
		<-srv.Notify(grpcserver.Stopped)
	} else {
		logger.Logger("Server failed to start.")
	}
}