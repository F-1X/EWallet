package app

import (
	"encoding/json"
	"ewallet/pkg/model"
	"log"
	"net/http"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, err.Error())
		}
	}
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	log.Println("some magic in writeJSon")
	w.Header().Add("Content-Type", "application/json")
	switch t := v.(type) {
	case model.ApiMessage:
		w.WriteHeader(t.ApiStatus)
		json.NewEncoder(w).Encode(model.ApiMessage{ApiMessage: t.String()})

	case model.ApiError:
		w.WriteHeader(t.ApiStatus)
		json.NewEncoder(w).Encode(model.ApiError{ApiError: t.String()})

	}
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(v)
}

func ConvertWalletToWalletResponce(wallet *model.Wallet) *model.WalletResponse {
	return &model.WalletResponse{
		ID:      wallet.WalletId,
		Balance: wallet.Balance,
	}
}

func ConvertTransactionToHistoryResponse(transactions []model.Transaction) []model.HistoryResponse {
	var historyResponses []model.HistoryResponse

	for _, transaction := range transactions {
		historyResponse := model.HistoryResponse{
			Time:   transaction.Time,
			From:   transaction.From,
			To:     transaction.To,
			Amount: transaction.Amount,
		}
		historyResponses = append(historyResponses, historyResponse)
	}

	return historyResponses
}
