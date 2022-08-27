package wirepod

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	leopard "github.com/Picovoice/leopard/binding/go"
	"github.com/asticode/go-asticoqui"
	"github.com/digital-dream-labs/chipper/pkg/logger"
	"github.com/pkg/errors"
)

var usePicovoice bool = false
var leopardSTTArray []leopard.Leopard
var picovoiceInstancesOS string
var picovoiceInstances int

const (
	// FallbackIntent is the failure-mode intent response
	FallbackIntent          = "intent_system_unsupported"
	IntentWeather           = "intent_weather"
	IntentWeatherExtend     = "intent_weather_extend"
	IntentNoLocation        = "intent_weather_unknownlocation"
	IntentNoDefaultLocation = "intent_weather_nodefaultlocation"

	IntentClockSetTimer                    = "intent_clock_settimer"
	IntentClockSetTimerExtend              = "intent_clock_settimer_extend"
	IntentNamesUsername                    = "intent_names_username"
	IntentNamesUsernameExtend              = "intent_names_username_extend"
	IntentPlaySpecific                     = "intent_play_specific"
	IntentPlaySpecificExtend               = "intent_play_specific_extend"
	IntentMessaqePlayMessage               = "intent_message_playmessage"
	IntentMessagePlayMessageExtend         = "intent_message_playmessage_extend"
	IntentMessageRecordMessage             = "intent_message_recordmessage"
	IntentMessageRecordMessageExtend       = "intent_message_recordmessage_extend"
	IntentGlobalStop                       = "intent_global_stop"
	IntentGlobalStopExtend                 = "intent_global_stop_extend"
	IntentGlobalDelete                     = "intent_global_delete"
	IntentGlobalDeleteExtend               = "intent_global_delete_extend"
	IntentPhotoTake                        = "intent_photo_take"
	IntentPhotoTakeExtend                  = "intent_photo_take_extend"
	IntentSystemDiscovery                  = "intent_system_discovery"
	IntentSystemDiscoveryExtend            = "intent_system_discovery_extend"
	IntentImperativeVolumeLevelExtend      = "intent_imperative_volumelevel_extend"
	IntentImperativeEyeColorSpecificExtend = "intent_imperative_eyecolor_specific_extend"
)

// Server stores the config
type Server struct{}

// New returns a new server
func New() (*Server, error) {
	// if os.Getenv("DEBUG_LOGGING") != "true" && os.Getenv("DEBUG_LOGGING") != "false" {
	// 	logger.Logger("No valid value for DEBUG_LOGGING, setting to true")
	// 	debugLogging = true
	// } else {
	// 	if os.Getenv("DEBUG_LOGGING") == "true" {
	// 		debugLogging = true
	// 	} else {
	// 		debugLogging = false
	// 	}
	// }
	InitHoundify()
	if os.Getenv("STT_SERVICE") == "leopard" {
		logger.Logger("Using Leopard")
		usePicovoice = true
		var picovoiceKey string
		picovoiceKeyOS := os.Getenv("PICOVOICE_APIKEY")
		leopardKeyOS := os.Getenv("LEOPARD_APIKEY")
		if picovoiceInstancesOS == "" {
			picovoiceInstances = 3
		} else {
			picovoiceInstancesToInt, err := strconv.Atoi(picovoiceInstancesOS)
			picovoiceInstances = picovoiceInstancesToInt
			if err != nil {
				fmt.Println("PICOVOICE_INSTANCES is not a valid integer, using default value of 3")
				picovoiceInstances = 3
			}
		}
		if picovoiceKeyOS == "" {
			if leopardKeyOS == "" {
				logger.Logger("You must set PICOVOICE_APIKEY to a value.")
				return nil, nil
			} else {
				fmt.Println("PICOVOICE_APIKEY is not set, using LEOPARD_APIKEY")
				picovoiceKey = leopardKeyOS
			}
		} else {
			picovoiceKey = picovoiceKeyOS
		}
		logger.Logger("Initializing " + strconv.Itoa(picovoiceInstances) + " Picovoice Instances...")
		for i := 0; i < picovoiceInstances; i++ {
			fmt.Println("Initializing Picovoice Instance " + strconv.Itoa(i))
			leopardSTTArray = append(leopardSTTArray, leopard.Leopard{AccessKey: picovoiceKey})
			leopardSTTArray[i].Init()
		}
	} else {
		usePicovoice = false
		var testTimer float64
		var timerDie bool = false
		logger.Logger("Running a Coqui test...")
		coquiInstance, _ := asticoqui.New("../stt/model.tflite")
		if _, err := os.Stat("../stt/large_vocabulary.scorer"); err == nil {
			coquiInstance.EnableExternalScorer("../stt/large_vocabulary.scorer")
		} else if _, err := os.Stat("../stt/model.scorer"); err == nil {
			err := coquiInstance.EnableExternalScorer("../stt/model.scorer")
			if err != nil {
				return nil, err
			}
		} else {
			logger.Logger("No .scorer file found.")
			return nil, errors.New("ERR: No .scorer file found.")
		}
		coquiStream, err := coquiInstance.NewStream()
		if err != nil {
			logger.Logger(err)
			return nil, nil
		}
		pcmBytes, _ := os.ReadFile("./stttest.pcm")
		var micData [][]byte
		micData = split(pcmBytes)
		for _, sample := range micData {
			coquiStream.FeedAudioContent(bytesToSamples(sample))
		}
		go func() {
			for testTimer <= 7.00 {
				if timerDie {
					break
				}
				time.Sleep(time.Millisecond * 10)
				testTimer = testTimer + 0.01
				if testTimer > 6.50 {
					logger.Logger("The STT test is taking too long, this hardware may not be adequate.")
				}
			}
		}()
		res, err := coquiStream.Finish()
		if err != nil {
			log.Fatal("Failed testing speech to text: ", err)
		}
		logger.Logger("Text:", res)
		timerDie = true
		logger.Logger("Coqui test successful! (Took " + strconv.FormatFloat(testTimer, 'f', 2, 64) + " seconds)")
	}
	return &Server{}, nil
}
