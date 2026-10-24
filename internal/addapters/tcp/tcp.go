package addapters

import (
	"encoding/json"
	"fmt"
	"net"
	"sockets-tc/config"
	"time"
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

	var payl Request
	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))

		err := dec.Decode(&payl)

		if err != nil {
			fmt.Println(err)
			return
		}

		handl, ok := MethodMap[payl.Method]
		if !ok {
			enc.Encode([]byte("invlid method"))
			return
		}

		handl(conn, payl)

	}

}
