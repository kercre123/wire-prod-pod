package sdkapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/digital-dream-labs/chipper/pkg/jdocsserver"
	"github.com/digital-dream-labs/chipper/pkg/logger"
	"github.com/fforchino/vector-go-sdk/pkg/vectorpb"
	"github.com/go-ping/ping"
)

// the big workaround

// the escape pod CA cert only gets appended to the cert store when a jdocs connection is created
// this doesn't happen at every boot
// this utilizes Vector's connCheck to see if a bot has disconnected from the server for more than 10 seconds
// if it has, it will pull jdocs from the bot which will cause the CA cert to get appended to the store

func pingJdocs(target string) {
	target = strings.Split(target, ":")[0]
	var serial string
	jsonBytes, err := os.ReadFile(jdocsserver.InfoPath)
	if err != nil {
		logger.Logger("Error opening " + jdocsserver.InfoPath + ", this bot likely hasn't been authed")
		logger.Logger("Error pinging jdocs")
		return
	}
	var robotSDKInfo RobotInfoStore
	json.Unmarshal(jsonBytes, &robotSDKInfo)
	matched := false
	for _, robot := range robotSDKInfo.Robots {
		if strings.TrimSpace(strings.ToLower(robot.IPAddress)) == strings.TrimSpace(strings.ToLower(target)) {
			matched = true
			serial = robot.Esn
		}
	}
	if !matched {
		logger.Logger("vector-go-sdk error: serial did not match any bot in bot json")
		logger.Logger("Error pinging jdocs")
		return
	}
	robotTmp, err := NewWP(serial, false)
	if err != nil {
		logger.Logger(err)
		return
	}
	sdkAddress = robotTmp.Cfg.Target
	robotGUID = robotTmp.Cfg.Token
	_, err = robotTmp.Conn.BatteryState(ctx, &vectorpb.BatteryStateRequest{})
	if err != nil {
		robotTmp, err = NewWP(serial, true)
		if err != nil {
			logger.Logger(err)
			logger.Logger("Error pinging jdocs")
			return
		}
		sdkAddress = robotTmp.Cfg.Target
		robotGUID = robotTmp.Cfg.Token
		_, err = robotTmp.Conn.BatteryState(ctx, &vectorpb.BatteryStateRequest{})
		if err != nil {
			logger.Logger("Error pinging jdocs, likely unauthenticated")
			return
		}
	}
	resp, err := robotTmp.Conn.PullJdocs(ctx, &vectorpb.PullJdocsRequest{
		JdocTypes: []vectorpb.JdocType{vectorpb.JdocType_ROBOT_SETTINGS},
	})
	if err != nil {
		logger.Logger(err)
		logger.Logger("Failed to pull jdocs")
		return
	}
	logger.Logger("Successfully got jdocs from " + serial)
	// write to file
	writeBytes, _ := json.Marshal(resp.NamedJdocs[0].Doc)
	os.WriteFile("./jdocs/vic:"+serial+"-vic.RobotSettings.json", writeBytes, 0644)
	return
}

var jdocsTargets []string
var jdocsTimers []int
var jdocsShouldPing []bool
var jdocsTimerStarted []bool
var jdocsTimerReset []bool

func startJdocsTimer(target string) {
	var jdocsBotNum int
	for num, ip := range jdocsTargets {
		if ip == target {
			jdocsBotNum = num
		}
	}
	if !jdocsTimerStarted[jdocsBotNum] {
		jdocsTimerStarted[jdocsBotNum] = true
		jdocsShouldPing[jdocsBotNum] = false
		logger.Logger("Starting jdocs pinger timer for " + target)
		go func() {
			// wait 10 seconds
			for {
				time.Sleep(time.Second * 1)
				jdocsTimers[jdocsBotNum] = jdocsTimers[jdocsBotNum] + 1
				if jdocsTimers[jdocsBotNum] == 10 {
					logger.Logger("No connCheck from " + target + " in more than 10 seconds, will ping jdocs on next check")
					jdocsShouldPing[jdocsBotNum] = true
					jdocsTimerStarted[jdocsBotNum] = false
					return
				}
				if jdocsTimerReset[jdocsBotNum] {
					jdocsTimers[jdocsBotNum] = 0
					//logger.Logger("Resetting timer to 0 for bot " + target)
					jdocsTimerReset[jdocsBotNum] = false
				}
			}
		}()
	}
}

func jdocsPingTimer(target string) bool {
	for num, ip := range jdocsTargets {
		if ip == target {
			var returnValue bool = jdocsShouldPing[num]
			startJdocsTimer(target)
			jdocsTimerReset[num] = true
			if returnValue {
				jdocsShouldPing[num] = false
			}
			return returnValue
		}
	}
	jdocsTargets = append(jdocsTargets, target)
	jdocsTimers = append(jdocsTimers, 0)
	jdocsShouldPing = append(jdocsShouldPing, false)
	jdocsTimerStarted = append(jdocsTimerStarted, false)
	jdocsTimerReset = append(jdocsTimerReset, false)
	startJdocsTimer(target)
	return true
}

func connCheck(w http.ResponseWriter, r *http.Request) {
	switch {
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	case r.URL.Path == "/ok:80":
		//	logger.Logger("connCheck request from " + r.RemoteAddr)
		robotTarget := strings.Split(r.RemoteAddr, ":")[0] + ":443"
		robotTargetCheck := strings.Split(r.RemoteAddr, ":")[0]
		jsonB, _ := os.ReadFile(jdocsserver.InfoPath)
		json := string(jsonB)
		if strings.Contains(json, strings.TrimSpace(robotTargetCheck)) {
			ping := jdocsPingTimer(robotTarget)
			if ping {
				pingJdocs(robotTarget)
			}
		}
		fmt.Fprintf(w, "ok")
		return
	case r.URL.Path == "/link-esn-and-target":
		esn := r.FormValue("esn")
		target := r.FormValue("target")
		logger.Logger(len([]rune(esn)))
		if len([]rune(esn)) != 8 {
			fmt.Fprintf(w, "failed to link bot: Serial number should equal 8 characters")
			return
		}
		pinger, err := ping.NewPinger(target)
		pinger.SetPrivileged(true)
		if err != nil {
			fmt.Fprintf(w, "failed to link bot: IP address not valid")
			return
		}
		pinger.Count = 1
		pinger.Timeout = time.Second * 2
		err = pinger.Run()
		if err != nil {
			logger.Logger(err)
			fmt.Fprintf(w, "failed to link bot: Couldn't ping bot, make sure you have entered the correct ip address")
			return
		}
		jdocsserver.StoreBotInfoStrings(target, esn)
		fmt.Fprintf(w, "success")
		return
	}
}
