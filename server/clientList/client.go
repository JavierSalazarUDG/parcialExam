package clientList

import "net"

type Client struct {
	Client   net.Conn
	Nickname string
}
