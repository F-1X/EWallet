package app

import (
	"github.com/gorilla/mux"
)

func (s *APIServer) NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(Logging)

	r.HandleFunc("/login", makeHTTPHandleFunc(loginHandler))
	r.HandleFunc("/logout", makeHTTPHandleFunc(logoutHandler)).Methods("POST")

	r.HandleFunc("/home", makeHTTPHandleFunc(homeHandler))

	r.HandleFunc("/register", makeHTTPHandleFunc(registerHandler))

	apiV1Router := r.PathPrefix("/api/v1/wallet").Subrouter()
	apiV1Router.Use(sessionHandler)
	apiV1Router.HandleFunc("", makeHTTPHandleFunc(s.handleCreateWallet)).Methods("POST")
	apiV1Router.HandleFunc("/{walletId}", makeHTTPHandleFunc(s.handleGetWallet)).Methods("GET")
	apiV1Router.HandleFunc("/{walletId}/send", makeHTTPHandleFunc(s.handleTranferWallet)).Methods("POST")
	apiV1Router.HandleFunc("/{walletId}/history", makeHTTPHandleFunc(s.handleHistory)).Methods("GET")

	return r
}
