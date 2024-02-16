package app

import (
	"ewallet/pkg/storage"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	storage    storage.Storage
	router     *mux.Router
}

func NewAPIServer(listenAddr string, storage storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage:    storage,
	}
}

func (s *APIServer) BindRouter(router *mux.Router) {
	s.router = router
}

func (s *APIServer) RunServer() {
	log.Println("[+] server started at", time.Now().Format("2006/01/02 15:04:05"))
	log.Fatal(http.ListenAndServe(s.listenAddr, s.router))
}
