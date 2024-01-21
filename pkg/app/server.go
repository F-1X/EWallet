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
	storage storage.Storage
}


func NewAPIServer(listenAddr string, storage storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		storage: storage,
	}
}


func (s *APIServer) RunServer() {

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/wallet",                    makeHTTPHandleFunc(s.handleCreateWallet)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}",         makeHTTPHandleFunc(s.handleGetWallet)).Methods("GET")
	router.HandleFunc("/api/v1/wallet/{walletId}/send",    makeHTTPHandleFunc(s.handleTranferWallet)).Methods("POST")
	router.HandleFunc("/api/v1/wallet/{walletId}/history", makeHTTPHandleFunc(s.handleHistory)).Methods("GET")


	log.Println("[+] server started at", time.Now().Format("2006-01-02 15:04:05"))
	log.Fatal(http.ListenAndServe(s.listenAddr, router))
}

