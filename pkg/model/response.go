package model

import (
	"errors"
	"net/http"
)

var (
	ErrInsufficientFuns                  = ApiError{"not enough money on wallet", http.StatusBadGateway}
	ErrIncorrectWalletOrTransactionError = ApiError{"Ошибка в пользовательском запросе или ошибка перевода", http.StatusBadRequest}
	ErrWalletNotFound                    = ApiError{"Исходящий кошелек не найден", http.StatusBadRequest}
	ErrIncorrectWallet                   = ApiError{"Некоректная запись кошелька", http.StatusBadRequest}
	ErrInvalidRequest                    = ApiError{"Ошибка в запросе", http.StatusBadRequest}
	ErrDuplicateWalletID                 = ApiError{"duplicate wallet", http.StatusBadRequest}
	ErrFromWalletNotFound                = ApiError{"Кошелек отправителя не найден", http.StatusBadRequest}
	ErrToWalletNotFound                  = ApiError{"Кошелек назначения не найден", http.StatusBadRequest}
	ErrHistoryEmpty                      = ApiError{"История отсутствует", http.StatusBadRequest}
	ErrLoadConfig                        = ApiError{"failed to load config", http.StatusBadRequest}
	TransactionSuccess                   = ApiMessage{"Перевод успешно проведен", http.StatusOK}
	WalletCreated                        = ApiMessage{"Кошелек создан", http.StatusOK}
)

type ApiError struct {
	ApiError  string `json:"error"`
	ApiStatus int
}

func (a *ApiError) Error() error {
	return errors.New(a.ApiError)
}

func (a *ApiError) String() string {
	return a.ApiError
}

func NewApiError(err string, status int) *ApiError {
	return &ApiError{
		ApiError:  err,
		ApiStatus: status,
	}
}

type ApiMessage struct {
	ApiMessage string `json:"message"`
	ApiStatus  int
}

func NewApiMessage(message string, status int) *ApiMessage {
	return &ApiMessage{
		ApiMessage: message,
		ApiStatus:  status,
	}
}

func (a *ApiMessage) String() string {
	return a.ApiMessage
}
