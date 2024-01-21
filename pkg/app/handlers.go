package app

import (
	"encoding/json"
	"ewallet/pkg/model"
	"ewallet/internal/util"
	"ewallet/internal/validator"
	"net/http"

	"github.com/gorilla/mux"
)


func (s *APIServer) handleCreateWallet(w http.ResponseWriter, r *http.Request) error {
	wallet_id := util.GenerateWallet()
	// проверка уникальности ID от коллизий.
	// возвращается 400 ошибка "Ошибка в запросе" 
	if err := s.storage.CheckUniqueWalletId(wallet_id); err != nil {
		if err == model.ErrDuplicateWalletID {
			for {
				wallet_id = util.GenerateWallet()
				if err := s.storage.CheckUniqueWalletId(wallet_id); err != model.ErrDuplicateWalletID {
					if err != nil {
						return WriteJSON(w,http.StatusBadRequest,model.ErrInvalidRequest)
					}
					break
				}
			}
		} else {
			return WriteJSON(w,http.StatusBadRequest,model.ErrInvalidRequest)
		}
	}
	walletId, err  := s.storage.CreateWallet(wallet_id)
	if err != nil {
		return WriteJSON(w,http.StatusBadRequest, model.ErrInvalidRequest)
	}

	return WriteJSON(w,http.StatusOK,ConvertWalletToWalletRequest(walletId))
}

func (s *APIServer) handleGetWallet(w http.ResponseWriter, r *http.Request) error {
	walletIdStr := mux.Vars(r)["walletId"]
	// валидация ID 
	if err := validator.ValidateWallet(walletIdStr); err != nil {
		return WriteJSON(w,http.StatusNotFound, err)
	}

	walletId, err := s.storage.GetWallet(walletIdStr)
	if err != nil {
		if err == model.ErrWalletNotFound{
			return WriteJSON(w,http.StatusNotFound,model.ErrWalletNotFound)
		}
		return WriteJSON(w,http.StatusInternalServerError,"internal error get wallet")
	}


	return WriteJSON(w, http.StatusOK, ConvertWalletToWalletRequest(walletId))
}


func (s *APIServer) handleTranferWallet(w http.ResponseWriter, r *http.Request) error {
	fromID := mux.Vars(r)["walletId"]

	// валидация ID 
	if err := validator.ValidateWallet(fromID); err != nil {
		return WriteJSON(w,http.StatusNotFound, err)
	}

	// валидация запроса
	req := new(model.TransactionRequest)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return WriteJSON(w, http.StatusBadRequest, model.ErrIncorrectWalletOrTransactionError)
	}

	// валидация ID 
	if err := validator.ValidateWallet(req.ToID); err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}

	// отрицательный перевод
	if req.Amount.IsNegative(){
		return WriteJSON(w, http.StatusBadRequest, model.ErrIncorrectWalletOrTransactionError)
	}
	
	_, err := s.storage.Transaction(fromID, req.ToID, req.Amount)
	if err != nil {
		if err == model.ErrIncorrectWalletOrTransactionError || err == model.ErrToWalletNotFound || err == model.ErrFromWalletNotFound {
			return WriteJSON(w,http.StatusInternalServerError, err)
		} else {
			// непредвиденные ошибки 
			return WriteJSON(w, http.StatusInternalServerError, model.ErrIncorrectWalletOrTransactionError)
		}
	}

	return WriteJSON(w, http.StatusOK, model.TransactionSuccess)
}

func (s *APIServer) handleHistory(w http.ResponseWriter, r *http.Request) error {
	walletIdStr := mux.Vars(r)["walletId"]
	// валидация ID 
	if err := validator.ValidateWallet(walletIdStr); err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}

	if _,err := s.storage.GetWallet(walletIdStr); err != nil {
		if err == model.ErrIncorrectWallet {
			return WriteJSON(w,http.StatusNotFound,model.ErrWalletNotFound)
		}
		// непредвиденные ошибки 
		return WriteJSON(w,http.StatusNotFound, "internal get wallet error")
		
	}
	
	tr,err := s.storage.History(walletIdStr)
	if err != nil {
		if err != model.ErrIncorrectWalletOrTransactionError || err != model.ErrHistoryEmpty{
			// непредвиденные ошибки 
			return WriteJSON(w,http.StatusInternalServerError,"internal history error")	
		}
		return WriteJSON(w,http.StatusInternalServerError,err)
		
	}

	return WriteJSON(w,http.StatusOK,ConvertTransactionToHistoryResponse(*tr))
}