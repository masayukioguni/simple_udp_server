package main

import (
	//"./bcd"
	"./payload"
	"./win"
	"fmt"
	"log"
	"net"
	"time"
)

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
		//log.Println("receivePayload currentPayload:", bufferLength)

		payloadChannel <- currentPayload
	}

	return nil
}

func processPayload(payloadChannel chan *payload.Payload) error {
	for {
		currentPayload := <-payloadChannel
		Winformat := &win.WinFormat{}
		Winformat = win.Parse(currentPayload.Buffer)
		fmt.Printf("%0x%0x\n", Winformat.Sequence, Winformat.SubSequence)
		//log.Println("processPayload currentPayload:", currentPayload)
		//log.Println("currentPayload.buffer:%04hX ", currentPayload.Buffer[0:1])
		/*
			seq := currentPayload.Buffer[0:2]
			A0 := currentPayload.Buffer[2:3]
			length := currentPayload.Buffer[3:5]
			bcddate := currentPayload.Buffer[5:11]

			year := bcd.BcdToInt(int(bcddate[0]))
			month := bcd.BcdToInt(int(bcddate[1]))
			day := bcd.BcdToInt(int(bcddate[2]))
			hour := bcd.BcdToInt(int(bcddate[3]))
			minute := bcd.BcdToInt(int(bcddate[4]))
			second := bcd.BcdToInt(int(bcddate[5]))

			ch := currentPayload.Buffer[11:13]
			size := currentPayload.Buffer[13] >> 4
			rate := (currentPayload.Buffer[13]&0x0f)<<8 | currentPayload.Buffer[14]&0xff

			firstSample := currentPayload.Buffer[15:19]

			//startPos := 19

			datetime := fmt.Sprintf("%02d%02d%02d%02d%02d%02d", year, month, day, hour, minute, second)
			fmt.Printf("%04X %X %X %s %04X %d %d %08X\n", seq, A0, length, datetime, ch, size, rate, firstSample)
		*/
		/*
			if size == 0 {
				rate = rate / 2

			}

			for i := 0; i < int(rate)-1; i++ {
				if size == 4 {
					s := int(startPos + (i * 4))
					e := int(startPos + (i * 4) + 4)

					diff := currentPayload.Buffer[s:e]
					fmt.Printf("s:%d e:%d index:%d %04X\n", s-19, e-19, i, diff)
				}
				if size == 3 {
					s := int(startPos + (i * 3))
					e := int(startPos + (i * 3) + 3)

					diff := currentPayload.Buffer[s:e]
					fmt.Printf("s:%d e:%d index:%d %04X\n", s-19, e-19, i, diff)
				}

				if size == 2 {
					s := int(startPos + (i * 2))
					e := int(startPos + (i * 2) + 2)

					diff := currentPayload.Buffer[s:e]
					fmt.Printf("s:%d e:%d index:%d %04X\n", s-19, e-19, i, diff)
				}

				if size == 1 {
					s := int(startPos + (i * 1))
					e := int(startPos + (i * 1) + 1)

					diff := currentPayload.Buffer[s:e]
					fmt.Printf("s:%d e:%d index:%d %X\n", s-19, e-19, i, diff)
				}

				if size == 0 {
					s := int(startPos + (i * 1))
					e := int(startPos + (i * 1) + 1)

					diff := currentPayload.Buffer[s:e]
					fmt.Printf("s:%d e:%d index:%d %X\n", s-19, e-19, i, diff)

				}
			}
		*/
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
