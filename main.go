package main

import (
	"github.com/gorilla/mux"
	"github.com/wdalmut/iot-http-to-low-level/board"
	"github.com/wdalmut/iot-http-to-low-level/proxy"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	server := proxy.Server{
		Router: router,
		HttpServer: &http.Server{
			Addr:    "0.0.0.0:8082",
			Handler: router,
		},
		BoardServer: &board.Server{
			Addr: "0.0.0.0:9005",
		},
	}

	server.ListenAndServe()
}
