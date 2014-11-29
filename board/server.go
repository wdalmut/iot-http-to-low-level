package board

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	Addr         string
	WriteChannel chan []byte
	ReadChannel  chan []byte
}

func (s *Server) ListenAndServe() {
	s.WriteChannel = make(chan []byte)
	s.ReadChannel = make(chan []byte)

	addr, _ := net.ResolveTCPAddr("tcp", s.Addr)
	ln, err := net.ListenTCP("tcp", addr)

	if err != nil {
		log.Panicf("Unable to bind port on %v", s.Addr)
	}

	for {
		conn, err := ln.AcceptTCP()
		conn.SetKeepAlive(false)

		if err == nil {
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	log.Println("New board connected")
	w := new(sync.WaitGroup)
	w.Add(1)
	go func() {
		for {
			data := <-s.WriteChannel

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

			if err != nil {
				conn.Close()
				break
			}

			if n > 0 {
				s.ReadChannel <- buffer
			}

		}

		defer w.Done()
	}()

	w.Wait()
	log.Println("Board Disconnected")
}
