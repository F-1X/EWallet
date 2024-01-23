package app

import (
	"encoding/json"
	"ewallet/pkg/model"
	"net/http"
)


type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w,r); err != nil {
			WriteJSON(w, http.StatusBadRequest, err.Error())
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type","application/json")
	switch t := v.(type) {
	case string:
		return json.NewEncoder(w).Encode(model.ApiMessage{Message:t})
	case error:
		return json.NewEncoder(w).Encode(model.ApiError{Error: t.Error()})
	}
	return json.NewEncoder(w).Encode(v)
}

func ConvertWalletToWalletRequest(wallet *model.Wallet) *model.WalletResponse {
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
