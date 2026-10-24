package addapters

import "net"

type HandlerFunc func(net.Conn, Request)

var MethodMap = map[Method]HandlerFunc{
	CreateEvent: CreateEventHandler,
}
