package main

import (
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

func main() {

}
