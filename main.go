package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	ListenAddr string
	listen     net.Listener
	quitCh     chan struct{}
	msgCh      chan Message
}

func NewServer(Addr string) *Server {
	return &Server{
		ListenAddr: Addr,
		quitCh:     make(chan struct{}),
		msgCh:      make(chan Message, 10),
	}
}

func (s *Server) StartServer() error {
	listen, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listen.Close()

	s.listen = listen

	go s.AcceptLoop()
	<-s.quitCh

	close(s.msgCh)

	return nil
}

func (s *Server) AcceptLoop() {
	for {
		conn, err := s.listen.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("new connection to tcp server:", conn.RemoteAddr())

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
		//identity of client and message sent by client to tcp server
		s.msgCh <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:cn],
		}
	}

}

func main() {
	server := NewServer(":6000")

	go func() {
		for msg := range server.msgCh {
			fmt.Printf("recieved message from connection(%s):%s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.StartServer())
}
