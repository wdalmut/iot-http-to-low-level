package proxy

import (
	"github.com/gorilla/mux"
	"github.com/wdalmut/iot-http-to-low-level/board"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"sync"
)

type ChunkWrapper interface {
	Wrap([]byte) []byte
}

type Server struct {
	Router      *mux.Router
	HttpServer  *http.Server
	BoardServer *board.Server
	Wrapper     ChunkWrapper
}

func (s *Server) HomeHandler(writer http.ResponseWriter, req *http.Request) {
	home, _ := ioutil.ReadFile("index.html")
	writer.Write(home)
}

func (s *Server) ReadHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html")

	closer := httputil.NewChunkedWriter(writer)

	//Start frame for browsers
	buf := make([]byte, 2000)
	closer.Write(buf)
	if f, ok := writer.(http.Flusher); ok {
		f.Flush()
	}

	for {
		data := <-s.BoardServer.ReadChannel

		_, e := closer.Write(s.Wrapper.Wrap(data))
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}

		if e != nil {
			closer.Close()
			break
		}
	}

	log.Println("HTTP Stream client disconnected")
}

func (s *Server) DataHandler(writer http.ResponseWriter, req *http.Request) {
	data := req.PostFormValue("data")

	s.BoardServer.WriteChannel <- []byte(data)

	writer.Write([]byte("OK"))
}

func (s *Server) ListenAndServe() {
	s.Router.HandleFunc("/board", s.HomeHandler).Methods("GET")
	s.Router.HandleFunc("/board/read", s.ReadHandler).Methods("GET")
	s.Router.HandleFunc("/board", s.DataHandler).Methods("POST")

	servers := new(sync.WaitGroup)

	servers.Add(1)
	go s.HttpServer.ListenAndServe()
	servers.Add(1)
	go s.BoardServer.ListenAndServe()

	servers.Wait()
}
