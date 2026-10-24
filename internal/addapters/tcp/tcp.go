package addapters

import (
	"encoding/json"
	"fmt"
	"net"
	"sockets-tc/config"
)

type Server struct {
	Env *config.Env
}

func NewServer(env *config.Env) *Server {
	return &Server{
		Env: env,
	}
}

func (s *Server) Start() error {

	listener, err := net.Listen("tcp", s.Env.Port)
	if err != nil {
		return err
	}

	defer listener.Close()
	fmt.Printf("Server listen on port %s ....  \n ", s.Env.Port)

	for {
		conn, err := listener.Accept()

		fmt.Println("connection accepted: ", conn)

		if err != nil {
			fmt.Println("failed to accept connection")
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	var payl Payload
	err := json.NewDecoder(conn).Decode(&payl)

	if err != nil {
		fmt.Println(err)
		return
	}

	switch payl.Method {
	case CreateTicket:
		conn.Write([]byte("ok"))
	default:
		conn.Write([]byte("unexpected method"))
	}

}
