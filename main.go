package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	ListenAddr string
	listen     net.Listener
	quitCh     chan struct{}
}

func NewServer(Addr string) *Server {
	return &Server{
		ListenAddr: Addr,
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) StartServer() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listen.Close()

	s.listen = listen
	<-s.quitCh

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		// each time there is a new connecton a new go routine is spun up to have a
		//large number of non blocking connections
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	//switch size of bytes depending on size of message being sent
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		cn, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}
		//stringify message sent by client to tcp server
		msg := buf[:cn]
		fmt.Println(string(msg))
	}

}

func main() {
	server := NewServer(":6000")
	log.Fatal(server.StartServer())
}
