package payload

import (
	"net"
)

type Payload struct {
	Addr         *net.UDPAddr
	Conn         *net.UDPConn
	Buffer       []byte
	BufferLength int
	Err          error
}
