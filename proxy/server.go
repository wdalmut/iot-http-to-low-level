package proxy

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/wdalmut/iot-http-to-low-level/board"
	"io/ioutil"
	"net/http"
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

func (s *Server) DataHandler(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)

	fmt.Printf("%v\n", vars)

	data := []byte(vars["data"])
	s.BoardServer.Write(data)

	writer.Write([]byte("sent"))

}

func (s *Server) ListenAndServe() {
	s.Router.HandleFunc("/board", s.HomeHandler).Methods("GET")
	s.Router.HandleFunc("/board/{data}", s.DataHandler).Methods("GET", "POST")

	servers := new(sync.WaitGroup)

	servers.Add(1)
	go s.HttpServer.ListenAndServe()
	servers.Add(1)
	go s.BoardServer.ListenAndServe()

	servers.Wait()
}
