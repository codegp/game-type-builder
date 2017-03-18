package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"teamrunner/api"
	"teamrunner/bot"
	"teamrunner/turninformer"
	"teamrunner/yielder"

	"git.apache.org/thrift.git/lib/go/thrift"
)

const retryLimit int = 120

var gameRunnerAddr string

const addr = 9000

func init() {
	gameRunnerAddr = "game-runner"
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
	yielders map[int32]*yielder.Yielder
}

func NewTurnInformerHandler(lifeOver chan bool, apiClient *api.APIClient) *TurnInformerHandler {
	ti := &TurnInformerHandler{
		lifeOver: lifeOver,
		api:      apiClient,
		yielders: map[int32]*yielder.Yielder{},
	}

	return ti
}

func (turninformer *TurnInformerHandler) CreateBot(botID int32) error {
	log.Println("IM ALIVE! ", botID)
	// TODO: recover from panics and invoke a timeout
	y := yielder.NewYielder()
	turninformer.yielders[botID] = y
	go turninformer.start(y)
	return nil
}

func (turninformer *TurnInformerHandler) DestroyBot(botID int32) error {
	log.Println("IM DEAD! ", botID)
	turninformer.lifeOver <- true
	return nil
}

func (turninformer *TurnInformerHandler) StartTurn(botID int32) error {
	log.Println("START TURN! ", botID)
	// TODO: recover from panics and invoke a timeout
	turninformer.yielders[botID].WaitForYield()
	return nil
}

func (turninformer *TurnInformerHandler) Destroy() error {
	log.Println("WE ALL DEAD!")
	turninformer.lifeOver <- true
	return nil
}

func (turninformer *TurnInformerHandler) start(y *yielder.Yielder) {
	y.WaitForStart()
	bot.Run(turninformer.api, y)
}

func newAPIClient() (*api.APIClient, error) {
	transport, err := thrift.NewTSocket(fmt.Sprintf("%s:9000", gameRunnerAddr))
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
		log.Printf("Could not connect to %s:9000. Retrying in 1s: %v. ", gameRunnerAddr, err)
		time.Sleep(time.Millisecond * time.Duration(1000))
		retryCount++
	}

	return nil, fmt.Errorf("Failed to start client at %s:9000, err:\n%v", gameRunnerAddr, err)
}
