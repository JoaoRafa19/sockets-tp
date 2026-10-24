package addapters

import "net"

type HandlerFunc func(net.Conn, Request) Response

var MethodMap = map[Method]HandlerFunc{
	CreateEvent: CreateEventHandler,
}
