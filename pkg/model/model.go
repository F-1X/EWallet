package model

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)



type WalletResponse struct {
	ID      string          `json:"id"`
	Balance decimal.Decimal `json:"balance"`
}

type Wallet struct {
	ID         string
	WalletId   string          `json:"id"`
	Balance    decimal.Decimal `json:"balance"`
	Created_at time.Time
	Updated_at time.Time
}

type Transaction struct {
	ID     string          `json:"id"`
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
	Time   time.Time       `json:"transfer_time"`
}

type HistoryResponse struct {
	Time   time.Time       `json:"time"`
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

type TransactionRequest struct {
	ToID   string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

func (hr *HistoryResponse) MarshalJSON() ([]byte, error) {
	amountFloat, _ := hr.Amount.Float64()
	return json.Marshal(struct {
		Time    time.Time `json:"id"`
		From    string    `json:"from"`
		To      string    `json:"to"`
		Balance float64   `json:"amount"`
	}{
		Time:    hr.Time,
		From:    hr.From,
		To:      hr.To,
		Balance: amountFloat,
	})
}

func (wr *WalletResponse) MarshalJSON() ([]byte, error) {
	balanceFloat, _ := wr.Balance.Float64()
	return json.Marshal(struct {
		ID      string  `json:"id"`
		Balance float64 `json:"balance"`
	}{
		ID:      wr.ID,
		Balance: balanceFloat,
	})
}

// Корректность передаваемых полех в теле запроса json хендлера перевода средств, транзакции
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
