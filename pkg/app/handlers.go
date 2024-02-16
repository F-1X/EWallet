package app

import (
	"encoding/json"
	"ewallet/internal/util"
	"ewallet/internal/validator"
	"ewallet/pkg/model"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func registerHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println("you are on register")
	if r.Method == "GET" {

		return WriteJSON(w, "register")
	}
	if r.Method == "POST" {
		return nil
	}
	return nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println("you are on home homeHandler")
	if r.Method == "GET" {

		return WriteJSON(w, "home")
	}

	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println("logouted ")
	cookie, err := r.Cookie(sessionKey)
	if err != nil {
		return nil
	}

	delete(sessions, cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:   sessionKey,
		MaxAge: -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)

	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) error {
	log.Println("you are on loggin")
	if r.Method == "GET" {
		return WriteJSON(w, "hello")
	}
	if r.Method == "POST" {
		var creds Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			return ErrDecodeCreds.Error()
		}
		password, ok := memoryCreds[creds.Username]
		if !ok {
			return ErrUserNotExist.Error()
		}
		if password != creds.Password {
			return ErrWrongPassword.Error()
		}

		sessindId := uuid.NewString()
		expires := time.Now().Add(5 * time.Minute) // 5 minutes
		maxAge := 1 * 60 * 60                      // 1 hour
		sessions[sessindId] = sessionCookie{
			Name:    sessionKey,
			Value:   sessindId,
			Expires: expires,
			MaxAge:  maxAge,
		}
		cookie := &http.Cookie{
			Name:    sessionKey,
			Value:   sessindId,
			Expires: expires,
			MaxAge:  maxAge,
		}
		http.SetCookie(w, cookie)
		log.Println("SetCookie", cookie)
		r.Header.Set("Location", "localhost:3000/home")
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return nil
	}

	return nil
}

func (s *APIServer) handleCreateWallet(w http.ResponseWriter, r *http.Request) error {
	wallet_id := util.GenerateWallet()
	if err := s.storage.CheckUniqueWalletId(wallet_id); err != nil {

		return model.ErrInvalidRequest.Error()
	}

	wallet, err := s.storage.CreateWallet(wallet_id)
	if err != nil {
		return model.ErrInvalidRequest.Error()
	}

	return WriteJSON(w, ConvertWalletToWalletResponce(wallet))
}

func (s *APIServer) handleGetWallet(w http.ResponseWriter, r *http.Request) error {
	log.Println("get wallet")
	walletIdStr := mux.Vars(r)["walletId"]
	if err := validator.ValidateWallet(walletIdStr); err != nil {

		return model.ErrWalletNotFound.Error()
	}

	wallet, err := s.storage.GetWallet(walletIdStr)
	if err != nil {
		return err
	}

	return WriteJSON(w, ConvertWalletToWalletResponce(wallet))
}

func (s *APIServer) handleTranferWallet(w http.ResponseWriter, r *http.Request) error {
	fromID := mux.Vars(r)["walletId"]

	if err := validator.ValidateWallet(fromID); err != nil {
		return model.ErrToWalletNotFound.Error()
	}

	req := new(model.TransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return model.ErrIncorrectWalletOrTransactionError.Error()
	}

	if err := validator.ValidateWallet(req.ToID); err != nil {
		return model.ErrToWalletNotFound.Error()
	}

	if req.Amount.IsNegative() {
		return model.ErrToWalletNotFound.Error()
	}

	_, err := s.storage.Transaction(fromID, req.ToID, req.Amount)
	if err != nil {
		return err
	}

	return WriteJSON(w, model.TransactionSuccess)
}

func (s *APIServer) handleHistory(w http.ResponseWriter, r *http.Request) error {

	walletIdStr := mux.Vars(r)["walletId"]
	if err := validator.ValidateWallet(walletIdStr); err != nil {
		return err
	}

	if _, err := s.storage.GetWallet(walletIdStr); err != nil {
		return err

	}

	tr, err := s.storage.History(walletIdStr)
	if err != nil {
		return err

	}

	return WriteJSON(w, ConvertTransactionToHistoryResponse(*tr))
}
