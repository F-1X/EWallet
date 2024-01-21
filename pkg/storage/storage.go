package storage

import (
	"ewallet/pkg/model"

	_ "github.com/lib/pq"
	"github.com/shopspring/decimal"
)

type Storage interface {
	CreateWallet(string) (*model.Wallet, error) 
	GetWallet(string) (*model.Wallet, error) 
	Transaction(string,string,decimal.Decimal) (*model.Transaction, error) 
	History(string) (*[]model.Transaction, error)
	CheckUniqueWalletId(string) error
}
