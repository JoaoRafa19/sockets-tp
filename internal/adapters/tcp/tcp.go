package tcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"sockets-tc/config"
)

type Server struct {
	Env     *config.Env
	handler *Handler
}

func NewServer(env *config.Env, handler *Handler) *Server {
	return &Server{
		Env:     env,
		handler: handler,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.Env.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Server listen on port %s ....\n", s.Env.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("failed to accept connection:", err)
			continue
		}

		fmt.Println("connection accepted:", conn.RemoteAddr())
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	dec := json.NewDecoder(conn)
	enc := json.NewEncoder(conn)

	for {
		conn.SetReadDeadline(time.Now().Add(10 * time.Minute))

		var payl Request
		if err := dec.Decode(&payl); err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) || isTimeout(err) {
				fmt.Println("connection closed:", conn.RemoteAddr())
				return
			}
			// JSON quebrado
			enc.Encode(fail("invalid payload: " + err.Error()))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp := s.handler.Handle(ctx, payl)
		cancel()

		if err := enc.Encode(resp); err != nil {
			fmt.Println("failed to write response:", err)
			return
		}
	}
}

func isTimeout(err error) bool {
	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}
