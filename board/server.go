package board

import (
	"log"
	"net"
	"sync"
	"time"
)

type Server struct {
	Addr         string
	writeChannel chan []byte
	ReadChannel  chan []byte
}

func (s *Server) ListenAndServe() {
	s.writeChannel = make(chan []byte)
	s.ReadChannel = make(chan []byte)

	addr, _ := net.ResolveTCPAddr("tcp", s.Addr)
	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Panicf("Unable to bind port on %v", s.Addr)
	}

	for {
		conn, err := ln.AcceptTCP()
		conn.SetKeepAlivePeriod(5 * time.Second)

		if err == nil {
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) Write(data []byte) {
	s.writeChannel <- data
}

func (s *Server) handleConnection(conn net.Conn) {
	w := new(sync.WaitGroup)
	w.Add(1)
	go func() {
		for {
			data := <-s.writeChannel

			_, err := conn.Write(data)
			if err != nil {
				conn.Close()
				break
			}
		}

		defer w.Done()
	}()

	w.Add(1)
	go func() {
		for {
			buffer := make([]byte, 128)
			n, err := conn.Read(buffer)

			if n > 0 {
				s.ReadChannel <- buffer
			}

			if err != nil {
				conn.Close()
			}
		}

		defer w.Done()
	}()
	w.Wait()

	log.Println("client disconnesso")
}