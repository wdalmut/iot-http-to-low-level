package board

import (
	"log"
	"net"
	"time"
)

type Server struct {
	Addr         string
	writeChannel chan byte
}

func (s *Server) ListenAndServe() {
	s.writeChannel = make(chan byte)

	addr, _ := net.ResolveTCPAddr("tcp", s.Addr)
	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Panicf("Unable to bind port on %v", s.Addr)
	}

	for {
		conn, err := ln.AcceptTCP()
		conn.SetKeepAlivePeriod(5 * time.Second)

		if err != nil {
			log.Printf("Non sono riuscito ad accettare il client %v\n", err)
		} else {
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) Write(data []byte) {
	s.writeChannel <- data[0]
}

func (s *Server) handleConnection(conn net.Conn) {
	var data byte

	for {
		data = <-s.writeChannel

		log.Println("%v", data)

		message := make([]byte, 1)
		message[0] = data
		_, err := conn.Write(message)

		if err != nil {
			conn.Close()
			break
		}
	}

	log.Println("client disconnesso")
}
