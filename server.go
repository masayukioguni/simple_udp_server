package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type payload struct {
	addr         *net.UDPAddr
	conn         *net.UDPConn
	buffer       []byte
	bufferLength int
	err          error
}

func loopToReadPayload(payloadChannel chan *payload,
	udpConn *net.UDPConn) error {
	for {
		buffer := make([]byte, 1400)

		bufferLength, udpAddr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Failed udpConn.ReadFromUDP():", err)
			continue
		}

		currentPayload := new(payload)
		currentPayload.addr = udpAddr
		currentPayload.conn = udpConn
		currentPayload.buffer = buffer
		currentPayload.bufferLength = bufferLength
		log.Println("loopToReadPayload currentPayload:", currentPayload)

		payloadChannel <- currentPayload
	}

	return nil
}

func loopToHandlePayload(payloadChannel chan *payload) error {
	for {
		currentPayload := <-payloadChannel
		log.Println("loopToHandlePayload currentPayload:", currentPayload)

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

	udpConn.SetReadBuffer(20000000)
	udpConn.SetWriteBuffer(20000000)

	payloadChannel := make(chan *payload)

	go loopToReadPayload(payloadChannel, udpConn)
	go loopToHandlePayload(payloadChannel)

	return nil
}

func main() {
	StartUdpServer(9229)
	time.Sleep(1000000000 * time.Second)
}
