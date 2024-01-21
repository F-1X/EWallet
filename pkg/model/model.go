package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)


var (
	
	ErrInsufficientFuns = errors.New("not enough money on wallet")

	ErrIncorrectWalletOrTransactionError = errors.New("Ошибка в пользовательском запросе или ошибка перевода")

	ErrWalletNotFound = errors.New("Исходящий кошелек не найден")

	ErrIncorrectWallet = errors.New("Некоректная запись кошелька")

	ErrInvalidRequest = errors.New("Ошибка в запросе")

	ErrDuplicateWalletID = errors.New("duplicate wallet")

	ErrFromWalletNotFound = errors.New("Кошелек отправителя не найден")
	ErrToWalletNotFound = errors.New("Кошелек назначения не найден")

	ErrHistoryEmpty = errors.New("История отсутствует")

	ErrLoadConfig = errors.New("failed to load config")

	TransactionSuccess = "Перевод успешно проведен"
	WalletCreated = "Кошелек создан"
)


type ApiError struct {
	Error  string            `json:"error"`
}

type ApiMessage struct {
	Message string           `json:"message"`
}

type WalletRequest struct {
	ID      string           `json:"id"`
	Balance decimal.Decimal  `json:"balance"`
}

type Wallet struct {
	ID       string            
	WalletId string			`json:"id"`
	Balance decimal.Decimal `json:"balance"`
	Created_at time.Time    
	Updated_at time.Time   
}

type Transaction struct {
	ID      string          `json:"id"`
	From    string          `json:"from"`
	To  	string          `json:"to"`
	Amount  decimal.Decimal `json:"amount"`
	Time    time.Time       `json:"transfer_time"`
}

type HistoryResponse struct {
	Time   time.Time         `json:"time"`
	From   string            `json:"from"`
	To     string            `json:"to"`
	Amount decimal.Decimal   `json:"amount"`
}

type TransactionRequest struct {
	ToID  	string          `json:"to"`
	Amount  decimal.Decimal `json:"amount"`
}


func (tr *TransactionRequest) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if to, ok := raw["to"]; ok {
		if err := json.Unmarshal(to, &tr.ToID); err != nil {
			return err
		}
	} else {
		return errors.New("invalidTO")
	}

	if amount, ok := raw["amount"]; ok {
		if err := json.Unmarshal(amount, &tr.Amount); err != nil {
			return err
		}
	} else {
		return errors.New("invalidAmount")
	}

	return nil
}