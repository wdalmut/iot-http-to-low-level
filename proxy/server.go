package proxy

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wdalmut/iot-http-to-low-level/board"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"sync"
)

type Server struct {
	Router      *mux.Router
	HttpServer  *http.Server
	BoardServer *board.Server
}

func (s *Server) HomeHandler(writer http.ResponseWriter, req *http.Request) {
	home, _ := ioutil.ReadFile("index.html")
	writer.Write(home)
}

func (s *Server) ReadHandler(writer http.ResponseWriter, req *http.Request) {
	closer := httputil.NewChunkedWriter(writer)

	//Start frame for browsers
	buf := make([]byte, 2000)
	closer.Write(buf)
	if f, ok := writer.(http.Flusher); ok {
		f.Flush()
	}

	for {
		data := <-s.BoardServer.ReadChannel

		_, e := closer.Write(data)
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}

		if e != nil {
			fmt.Printf("%v", e)
			break
		}
	}
}

func (s *Server) DataHandler(writer http.ResponseWriter, req *http.Request) {
	data := req.PostFormValue("data")
	s.BoardServer.Write([]byte(data))
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
