package app

import (
	"ewallet/pkg/model"
	"log"
	"net/http"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
		next.ServeHTTP(w, req)
	})
}

func sessionHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("sessionHandler___")
		if !isAuthenticated(r) {
			log.Println("Unathorized!")
			r.Header.Set("Location", "localhost:3000/login")
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		log.Println("Authorized!")
		next.ServeHTTP(w, r)
	})
}

func isAuthenticated(r *http.Request) bool {
	log.Println("checking auth")

	cookie, err := r.Cookie(sessionKey)
	if err != nil {
		return false
	}
	session, ok := sessions[cookie.Value]
	if !ok {
		return false
	}

	if time.Now().After(session.Expires) {
		return false
	}
	return true
}

var (
	ErrDecodeCreds   = model.ApiError{ApiError: "err decode creds", ApiStatus: http.StatusUnauthorized}
	ErrUserNotExist  = model.ApiError{ApiError: "not exist user", ApiStatus: http.StatusUnauthorized}
	ErrWrongPassword = model.ApiError{ApiError: "wrong password", ApiStatus: http.StatusUnauthorized}
)
