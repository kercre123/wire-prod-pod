package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/kercre123/chipper/pkg/initwirepod"
	"github.com/kercre123/chipper/pkg/logger"
	wpweb "github.com/kercre123/chipper/pkg/wirepod/config-ws"
	sdkWeb "github.com/kercre123/chipper/pkg/wirepod/sdkapp"
	coqui "github.com/kercre123/chipper/pkg/wirepod/stt/coqui"
	leopard "github.com/kercre123/chipper/pkg/wirepod/stt/leopard"
	vosk "github.com/kercre123/chipper/pkg/wirepod/stt/vosk"

	//	grpclog "github.com/digital-dream-labs/hugh/grpc/interceptors/log"

	grpcserver "github.com/digital-dream-labs/hugh/grpc/server"
	"github.com/digital-dream-labs/hugh/log"
)

var srv *grpcserver.Server

type chipperConfigStruct struct {
	Port                 string `json:"port"`
	Cert                 string `json:"cert"`
	Key                  string `json:"key"`
	WeatherEnable        string `json:"weatherEnable"`
	WeatherKey           string `json:"weatherKey"`
	WeatherUnit          string `json:"weatherUnit"`
	KnowledgeEnable      string `json:"knowledgeEnable"`
	KnowledgeProvider    string `json:"knowledgeProvider"`
	KnowledgeID          string `json:"knowledgeID"`
	KnowledgeKey         string `json:"knowledgeKey"`
	KnowledgeIntentGraph string `json:"knowledgeGraph"`
	// Houndify* is deprecated
	HoundifyEnable string `json:"houndifyEnable"`
	HoundifyKey    string `json:"houndifyKey"`
	HoundifyID     string `json:"houndifyID"`
	SttService     string `json:"sttService"`
	PicovoiceKey   string `json:"picovoiceKey"`
}

func chipperAPIHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case r.URL.Path == "/chipper/start_chipper":
		//	name := r.FormValue("name")
		if initwirepod.ServerRunning {
			fmt.Fprintf(w, "chipper already running")
		} else {
			startServer()
			fmt.Fprintf(w, "chipper started")
		}
		return
	case r.URL.Path == "/chipper/stop_chipper":
		if !initwirepod.ServerRunning {
			fmt.Fprintf(w, "chipper already stopped")
		} else {
			stopServer()
			fmt.Fprintf(w, "chipper stopped")
		}
		return
	case r.URL.Path == "/chipper/restart_chipper":
		if !initwirepod.ServerRunning {
			startServer()
			fmt.Fprintf(w, "chipper restarted")
		} else {
			stopServer()
			startServer()
			fmt.Fprintf(w, "chipper restarted")
		}
		return
	case r.URL.Path == "/chipper/make_config":
		port := r.FormValue("port")
		certType := r.FormValue("certType")
		weatherEnable := r.FormValue("weatherEnable")
		weatherKey := r.FormValue("weatherKey")
		weatherUnit := r.FormValue("weatherUnit")
		knowledgeEnable := r.FormValue("knowledgeEnable")
		knowledgeProvider := r.FormValue("knowledgeProvider")
		knowledgeID := r.FormValue("knowledgeID")
		knowledgeKey := r.FormValue("knowledgeKey")
		knowledgeIntent := r.FormValue("knowledgeIntent")
		// begin deprecation
		houndifyEnable := r.FormValue("houndifyEnable")
		// end deprecation
		sttService := r.FormValue("sttService")
		picovoiceKey := r.FormValue("picovoiceKey")
		var chipperConfigReq chipperConfigStruct
		chipperConfigReq.Port = port
		if strings.Contains(certType, "epod") {
			logger.Println("creating useepod")
			os.WriteFile("./useepod", []byte("true"), 0644)
			exec.Command("/bin/bash", "../setup.sh", "certs", "epod").Run()
			chipperConfigReq.Cert = "./epod/ep.crt"
			chipperConfigReq.Key = "./epod/ep.key"
		} else {
			exec.Command("/bin/rm", "-f", "./useepod").Run()
			cmdOutput, _ := exec.Command("/bin/bash", "../setup.sh", "certs", "ip").Output()
			if strings.Contains(string(cmdOutput), "Generating key and cert") {
				logger.Println("Successfully generated certs")
			}
			chipperConfigReq.Cert = "../certs/cert.crt"
			chipperConfigReq.Key = "../certs/cert.key"
		}
		if houndifyEnable != "" {
			logger.Println("houndifyEnable found in make config request, erroring")
			fmt.Fprintf(w, "failed: Your version of the webapp is too old, refresh it with CTRL + SHIFT + R and try again")
			return
		}
		chipperConfigReq.WeatherEnable = weatherEnable
		chipperConfigReq.WeatherKey = weatherKey
		chipperConfigReq.WeatherUnit = weatherUnit
		chipperConfigReq.KnowledgeEnable = knowledgeEnable
		chipperConfigReq.KnowledgeProvider = knowledgeProvider
		chipperConfigReq.KnowledgeKey = knowledgeKey
		chipperConfigReq.KnowledgeID = knowledgeID
		chipperConfigReq.SttService = sttService
		chipperConfigReq.PicovoiceKey = picovoiceKey
		chipperConfigReq.KnowledgeIntentGraph = knowledgeIntent
		chipperConfigBytes, _ := json.Marshal(chipperConfigReq)
		os.WriteFile("./chipperConfig.json", chipperConfigBytes, 0644)
		fmt.Fprintf(w, "chipper config created")
		return
	case r.URL.Path == "/chipper/upload_ssh_key":
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
		var serverConfig string
		if _, err := os.Stat("./useepod"); err == nil {
			serverConfig = `{"jdocs": "escapepod.local:443", "tms": "escapepod.local:443", "chipper": "escapepod.local:443", "check": "escapepod.local/ok:80", "logfiles": "s3://anki-device-logs-prod/victor", "appkey": "oDoa0quieSeir6goowai7f"}`
		} else {
			ipAddrBytes, err := os.ReadFile("../certs/address")
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}
			ipAddr := strings.TrimSpace(string(ipAddrBytes))
			serverConfig = `{"jdocs": "` + ipAddr + `:` + os.Getenv("DDL_RPC_PORT") + `", "tms": "` + ipAddr + `:` + os.Getenv("DDL_RPC_PORT") + `", "chipper": "` + ipAddr + `:` + os.Getenv("DDL_RPC_PORT") + `", "check": "` + ipAddr + `/ok:80` + `", "logfiles": "s3://anki-device-logs-prod/victor", "appkey": "oDoa0quieSeir6goowai7f"}`
		}
		exec.Command("/bin/mkdir", "-p", "../certs").Run()
		os.WriteFile("../certs/server_config.json", []byte(serverConfig), 0644)
		if _, err := os.Stat("/tmp/sshKey"); err != nil {
			cmdOutput, _ := exec.Command("/bin/bash", "./setupBot.sh", botIP).Output()
			fmt.Fprintf(w, "Output: "+string(cmdOutput))
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
			logger.Println("You must use the webserver to define where your certs are.")
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
		logger.Println("Use the webserver to setup and start chipper.")
	}
	var webPort string
	http.HandleFunc("/api/", wpweb.ApiHandler)
	go sdkWeb.BeginServer()
	http.HandleFunc("/chipper/", chipperAPIHandler)
	fileServer := http.FileServer(http.Dir("./webroot"))
	http.Handle("/", fileServer)
	if os.Getenv("WEBSERVER_PORT") != "" {
		if _, err := strconv.Atoi(os.Getenv("WEBSERVER_PORT")); err == nil {
			webPort = os.Getenv("WEBSERVER_PORT")
		} else {
			logger.Println("WEBSERVER_PORT contains letters, using default of 8080")
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
	initwirepod.StopGrpc()
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
		if chipperConfig.Cert == "./epod/ep.crt" {
			exec.Command("/bin/touch", "./useepod").Run()
		}
		os.Setenv("DDL_RPC_PORT", chipperConfig.Port)
		os.Setenv("DDL_RPC_TLS_CERTIFICATE", certString)
		os.Setenv("DDL_RPC_TLS_KEY", keyString)
		os.Setenv("DDL_RPC_CLIENT_AUTHENTICATION", "NoClientCert")
		os.Setenv("WEATHERAPI_ENABLED", chipperConfig.WeatherEnable)
		os.Setenv("WEATHERAPI_KEY", chipperConfig.WeatherKey)
		os.Setenv("WEATHERAPI_PROVIDER", "openweathermap.org")
		os.Setenv("WEATHERAPI_UNIT", chipperConfig.WeatherUnit)
		if chipperConfig.HoundifyEnable != "" {
			logger.Println("Old config version found, updating")
			os.Setenv("KNOWLEDGE_ENABLED", chipperConfig.HoundifyEnable)
			chipperConfig.KnowledgeEnable = chipperConfig.HoundifyEnable
			chipperConfig.KnowledgeIntentGraph = "false"
			if chipperConfig.HoundifyEnable == "true" {
				os.Setenv("KNOWLEDGE_PROVIDER", "houndify")
				os.Setenv("KNOWLEDGE_KEY", chipperConfig.HoundifyKey)
				os.Setenv("KNOWLEDGE_ID", chipperConfig.HoundifyID)
				chipperConfig.KnowledgeProvider = "houndify"
				chipperConfig.KnowledgeIntentGraph = "false"
				chipperConfig.KnowledgeID = chipperConfig.HoundifyID
				chipperConfig.KnowledgeKey = chipperConfig.HoundifyKey
				chipperConfig.HoundifyID = ""
				chipperConfig.HoundifyKey = ""
				chipperConfig.HoundifyEnable = ""
			}
			bytes, err := json.Marshal(chipperConfig)
			logger.Println("Updated json: " + string(bytes))
			if err != nil {
				logger.Println(err)
			}
			os.WriteFile("./chipperConfig.json", bytes, 0644)
		} else {
			os.Setenv("KNOWLEDGE_ENABLED", chipperConfig.KnowledgeEnable)
			os.Setenv("KNOWLEDGE_INTENT_GRAPH", chipperConfig.KnowledgeIntentGraph)
			os.Setenv("KNOWLEDGE_PROVIDER", chipperConfig.KnowledgeProvider)
			os.Setenv("KNOWLEDGE_KEY", chipperConfig.KnowledgeKey)
			os.Setenv("KNOWLEDGE_ID", chipperConfig.KnowledgeID)
		}
		os.Setenv("DEBUG_LOGGING", "true")
		os.Setenv("STT_SERVICE", chipperConfig.SttService)
		os.Setenv("PICOVOICE_APIKEY", chipperConfig.PicovoiceKey)
	}
	if !initwirepod.ServerRunning {
		if os.Getenv("STT_SERVICE") == "leopard" {
			initwirepod.StartServer(leopard.Init, leopard.STT, leopard.Name)
		} else if os.Getenv("STT_SERVICE") == "coqui" {
			initwirepod.StartServer(coqui.Init, coqui.STT, coqui.Name)
		} else {
			initwirepod.StartServer(vosk.Init, vosk.STT, vosk.Name)
		}
	}
}
