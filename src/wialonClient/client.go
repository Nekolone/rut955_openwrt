package wialonClient

import (
	"net"
	"sync"
)

type
struct {

}
type WialonClientHandler

type client struct {

	addr
	connection *net.TCPConn
	netWorkStatus string
}
type Client interface {
	ConntectToServer(addres string)
}

func NewWialonClient(handler WialonClientHandler) Client {
	return &client{}
}
