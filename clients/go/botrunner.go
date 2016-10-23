package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"botrunner/api"
	"botrunner/bot"
	"botrunner/turninformer"
	"botrunner/yielder"
	"git.apache.org/thrift.git/lib/go/thrift"
)

const retryLimit int = 20

var gameRunnerIP string

const addr = 9000

func init() {
	gameRunnerIP = os.Getenv("GAME_RUNNER_IP")
}

func main() {
	log.Printf("Starting informer server on %d\n", addr)
	lifeOver := make(chan bool)

	apiClient, err := newAPIClient()
	if err != nil {
		log.Fatalf("ERROR connecting to api server, %v", err)
	}

	server, err := getServer(lifeOver, apiClient)
	if err != nil {
		log.Fatalf("Failed to start team informer server %v", err)
	}

	go func() {
		log.Printf("Starting server")
		e := server.Serve()
		if e != nil {
			log.Fatalf("the errz!!! %v", e)
		}
	}()

	// botInformerClient, err := newBotInformerClient()
	// if err != nil {
	// 	log.Fatalf("Failed to start botInformerClient %v", err)
	// }
	//
	// botID := os.Getenv("BOT_ID")
	// podIP := os.Getenv("POD_IP")
	// if botID == "" || podIP == "" {
	// 	log.Fatalf("BotID and PodIP must be in env:\nbotID: %s\npodIP: %s", bodID, podIP)
	// }

	// botInformerClient.Started(botID, podIP)

	log.Println("team runner waiting for game to finish...")
	<-lifeOver
	server.Stop()
	log.Println("closing team runner")
}

func getServer(lifeOver chan bool, apiClient *api.APIClient) (*thrift.TSimpleServer, error) {
	//transport is a thrift.TServerTransport
	transport, err := thrift.NewTServerSocket(fmt.Sprintf("%s:%d", os.Getenv("POD_IP"), addr))
	if err != nil {
		return nil, err
	}
	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	handler := NewTurnInformerHandler(lifeOver, apiClient)
	processor := turninformer.NewTurnInformerProcessor(handler)
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)

	log.Printf("Starting the server on %d\n", addr)
	return server, nil
}

type TurnInformerHandler struct {
	lifeOver chan bool
	api      *api.APIClient
	yielder  *yielder.Yielder
}

func NewTurnInformerHandler(lifeOver chan bool, apiClient *api.APIClient) *TurnInformerHandler {
	ti := &TurnInformerHandler{
		lifeOver: lifeOver,
		api:      apiClient,
		yielder:  yielder.NewYielder(),
	}

	go ti.start()
	return ti
}

func (turninformer *TurnInformerHandler) StartTurn() error {
	log.Println("START TURN!")
	// TODO: recover from panics and invoke a timeout
	turninformer.yielder.WaitForYield()
	return nil
}

func (turninformer *TurnInformerHandler) Destroy() error {
	log.Println("IM DEAD!")
	turninformer.lifeOver <- true
	return nil
}

func (turninformer *TurnInformerHandler) start() {
	turninformer.yielder.WaitForStart()
	bot.Run(turninformer.api, turninformer.yielder)
	turninformer.lifeOver <- true
}

// func newBotInformerClient() (*botstartinformer.BotStartInformerClient, error) {
// 	transport, err := thrift.NewTSocket(fmt.Sprintf("http://%s:9000", gameRunnerIP))
// 	if err != nil {
// 		log.Printf("Error opening socket: %v\n", err)
// 		return nil, err
// 	}
//
// 	transportFactory := thrift.NewTTransportFactory()
// 	protocolFactory := thrift.NewTCompactProtocolFactory()
// 	t := transportFactory.GetTransport(transport)
//
// 	retryCount := 0
// 	for retryCount <= retryLimit {
// 		if err := t.Open(); err == nil {
// 			return botstartinformer.NewBotStartInformerClientFactory(t, protocolFactory), nil
// 		}
// 		time.Sleep(time.Millisecond * 100 * retryCount)
// 		retryCount++
// 	}
//
// 	return nil, fmt.Errorf("Failed to start client at port 9000")
// }

func newAPIClient() (*api.APIClient, error) {
	transport, err := thrift.NewTSocket(fmt.Sprintf("%s:9000", gameRunnerIP))
	if err != nil {
		log.Printf("Error opening socket: %v\n", err)
		return nil, err
	}

	transportFactory := thrift.NewTTransportFactory()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	t := transportFactory.GetTransport(transport)

	retryCount := 0
	for retryCount <= retryLimit {
		if err = t.Open(); err == nil {
			return api.NewAPIClientFactory(t, protocolFactory), nil
		}
		time.Sleep(time.Millisecond * time.Duration(100*retryCount))
		retryCount++
	}

	return nil, fmt.Errorf("Failed to start client at %s:9000, err:\n%v", gameRunnerIP, err)
}
