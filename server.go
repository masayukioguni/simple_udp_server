package main

import (
	"./payload"
	"encoding/json"
	"fmt"
	"github.com/masayukioguni/winformat"
	"github.com/t-k/fluent-logger-golang/fluent"
	"log"
	"net"
	"sync"
	"time"
)

const (
	defaultPort       = 9229
	defaultFluentHost = "127.0.0.1"
	defaultFluentPort = 24224
	defaultWindowSize = 1 * 1024 * 1024
	defaultTagName    = "win.format"
)

type Config struct {
	Port       int
	FluentHost string
	FluentPort int
	WindowSize int
	TagName    string
}

type Server struct {
	Config  Config
	conn    net.Conn
	pending []byte
	mutex   sync.Mutex
}

func New(config Config) (s *Server, err error) {
	if config.Port == 0 {
		config.Port = defaultPort
	}
	if config.FluentHost == "" {
		config.FluentHost = defaultFluentHost
	}
	if config.FluentPort == 0 {
		config.FluentPort = defaultFluentPort
	}
	if config.TagName == "" {
		config.TagName = defaultTagName
	}

	if config.WindowSize == 0 {
		config.WindowSize = defaultWindowSize
	}
	f = &Server{Config: config}

}

func receivePayloadProcess(payloadChannel chan *payload.Payload,
	udpConn *net.UDPConn) error {
	for {
		buffer := make([]byte, 1400)

		bufferLength, udpAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed udpConn.ReadFromUDP():", err)
			continue
		}

		currentPayload := new(payload.Payload)
		currentPayload.Addr = udpAddr
		currentPayload.Conn = udpConn
		currentPayload.Buffer = buffer
		currentPayload.BufferLength = bufferLength

		payloadChannel <- currentPayload
	}

	return nil
}

func processPayload(payloadChannel chan *payload.Payload) error {
	logger, err := fluent.New(fluent.Config{FluentPort: 24224, FluentHost: "127.0.0.1"})
	if err != nil {
		fmt.Println(err)
	}
	defer logger.Close()
	tag := "debug.access"

	for {
		currentPayload := <-payloadChannel
		jsonData, _ := json.Marshal(winformat.NewWinFormat(currentPayload.Buffer))
		logger.Post(tag, jsonData)
	}
	return nil
}

func StartUdpServer(udpPort int) error {
	log.Println("Trying to start UDP server port:", udpPort)

	udpServerAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", udpPort))
	if err != nil {
		log.Println("ResolveUDPAddr:", err)
		return err
	}
	log.Println("net.ResolveUDPAddr:", udpServerAddr)

	udpConn, err := net.ListenUDP("udp", udpServerAddr)
	if err != nil {
		return err
	}

	log.Println("net.ListenUDP:", udpConn)

	const bufferSize int = 1 * 1024 * 1024
	udpConn.SetReadBuffer(bufferSize)
	udpConn.SetWriteBuffer(bufferSize)

	payloadChannel := make(chan *payload.Payload)

	go receivePayloadProcess(payloadChannel, udpConn)
	go processPayload(payloadChannel)

	return nil
}

func main() {
	StartUdpServer(9229)
	time.Sleep(1000000000 * time.Second)
}
