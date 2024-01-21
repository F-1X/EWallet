package storage

import (
	"database/sql"
	"ewallet/pkg/model"
	"log"
	"time"

	"github.com/shopspring/decimal"
)

type Postgres struct {
	DB *sql.DB
}

func NewPostgres(DB_URL string) (*Postgres, error) {
	db, err := sql.Open("postgres", DB_URL)
	if err != nil {
		log.Fatal("failed open postgres db: ",err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("failed to connect to postgres: ",err)
	}
	return &Postgres{DB:db}, nil
}

func (s *Postgres) InitTables() error {
	if err := s.createWalletTable(); err != nil {
		return err
	}
	if err := s.createTransactionTable(); err != nil {
		return err
	}
	return nil
}

func (s *Postgres) createWalletTable() error {
	query := `CREATE TABLE IF NOT EXISTS
		wallet (
			id SERIAL NOT NULL PRIMARY KEY,
			wallet_id VARCHAR(64) NOT NULL UNIQUE, 
			balance NUMERIC NOT NULL,
			created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`
	_,err := s.DB.Exec(query)
	return err
}

func (s *Postgres) createTransactionTable() error {
	query := `CREATE TABLE IF NOT EXISTS
		transactions (
			id SERIAL PRIMARY KEY,
			from_wallet_id VARCHAR(64) NOT NULL,
			to_wallet_id VARCHAR(64) NOT NULL,
			amount NUMERIC NOT NULL,
			transfer_time TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`
	_,err := s.DB.Exec(query)
	return err
}

func (s *Postgres) CreateWallet(wallet_id string) (*model.Wallet, error) {
	log.Println("posgtres.CreateWallet wallet_id: ",wallet_id, " start")
	defer log.Println("posgtres.CreateWallet wallet_id: ",wallet_id," finished")

	query := "INSERT INTO wallet (wallet_id, balance) VALUES ($1, $2) RETURNING wallet_id, balance"

	stmt,err := s.DB.Prepare(query)
	if err != nil {
		log.Fatal(err,stmt)
	}
	defer stmt.Close()

	var walletId string
	var balance decimal.Decimal
	err = stmt.QueryRow(wallet_id, 100.0).Scan(&walletId, &balance)
	if err != nil {
		log.Fatal(err)
	}
	
	return &model.Wallet{WalletId: walletId, Balance: balance}, nil
}



func (s *Postgres) CheckUniqueWalletId(wallet_id string) error {
	log.Println("posgtres.CheckUniqueWalletId wallet_id: ",wallet_id, " started")
	defer log.Println("posgtres.CheckUniqueWalletId wallet_id: ",wallet_id, " finished")

	var existingWalletId string
	err := s.DB.QueryRow("SELECT wallet_id FROM wallet WHERE wallet_id = $1 LIMIT 1", wallet_id).Scan(&existingWalletId)
	if err != nil && err != sql.ErrNoRows{
		return err
	}
	
	if existingWalletId != "" {
        return model.ErrDuplicateWalletID
    }
	
	return nil
}

func (s *Postgres) GetWallet(wallet_id string) (*model.Wallet, error) {
	log.Println("posgtres.GetWallet wallet_id: ",wallet_id, " started")
	defer log.Println("posgtres.GetWallet wallet_id: ",wallet_id, " finished")

	rows, err := s.DB.Query("SELECT * FROM wallet WHERE wallet_id = $1", wallet_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		return scanIntoWallet(rows)
	}
	
	return nil , model.ErrIncorrectWallet
}

func (s *Postgres) Transaction(fromID string, toID string, amount decimal.Decimal) (*model.Transaction, error)  {
	log.Println("posgtres.Transaction fromID: ",fromID, " toID:",toID, "amount: ", amount, " started")
	defer log.Println("posgtres.Transaction fromID: ",fromID, " toID:",toID, "amount: ", amount, " finished")
	
	tx, err := s.DB.Begin()
	if err != nil {
		return nil,err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	var fromBalance decimal.Decimal
	err = tx.QueryRow("SELECT balance FROM wallet WHERE wallet_id = $1", fromID).Scan(&fromBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrFromWalletNotFound
		}
		tx.Rollback()
		return nil,err
	}
	if fromBalance.LessThan(amount) {
		log.Println("Balance " + fromBalance.String() + " less then amount " + amount.String() + " fromID: " + fromID)
		tx.Rollback()
		return nil, model.ErrIncorrectWalletOrTransactionError
	}

	fromBalance = fromBalance.Sub(amount)
	_, err = tx.Exec("UPDATE wallet SET balance = $1 WHERE wallet_id = $2", fromBalance, fromID)
	if err != nil {
		tx.Rollback()
		return nil,err
	}

	var toBalance decimal.Decimal
	err = tx.QueryRow("SELECT balance FROM wallet WHERE wallet_id = $1", toID).Scan(&toBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrToWalletNotFound
		}
		tx.Rollback()
		return nil,err
	}
	toBalance = toBalance.Add(amount)
	_, err = tx.Exec("UPDATE wallet SET balance = $1 WHERE wallet_id = $2", toBalance, toID)
	if err != nil {
		tx.Rollback()
		return nil,err
	}

	_, err = tx.Exec(`
		INSERT INTO transactions (from_wallet_id, to_wallet_id, amount, transfer_time)
		VALUES ($1, $2, $3, $4)
		`, fromID, toID, amount,time.Now().Format(time.RFC3339))
	if err != nil {
		tx.Rollback()
		return nil,err
	}

	err = tx.Commit()
	if err != nil {
		return nil,err
	}
	
	return nil,nil

}

func (s *Postgres) History(fromID string) (*[]model.Transaction, error)  {
	log.Println("posgtres.History fromID: ",fromID, " started")
	defer log.Println("posgtres.History fromID: ",fromID, " finished")

	outgoingQuery := `
		SELECT id, from_wallet_id, to_wallet_id, amount, transfer_time
		FROM transactions
		WHERE from_wallet_id = $1
	`

	incomingQuery := `
		SELECT id, from_wallet_id, to_wallet_id, amount, transfer_time
		FROM transactions
		WHERE to_wallet_id = $1
	`

	var transactions []model.Transaction

	rows, err := s.DB.Query(outgoingQuery, fromID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions, err = scanIntoTransaction(rows,transactions)
	if err != nil {
		return nil, err
	}

	rows, err = s.DB.Query(incomingQuery, fromID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	transactions, err = scanIntoTransaction(rows,transactions)
	if err != nil {
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, model.ErrHistoryEmpty
	}
	
	return &transactions, nil
}


func scanIntoWallet(rows *sql.Rows) (*model.Wallet, error) {
	wallet := new(model.Wallet)
	err := rows.Scan(
		&wallet.ID,
		&wallet.WalletId,
		&wallet.Balance,
		&wallet.Created_at,
		&wallet.Updated_at)

	return wallet, err
}

func scanIntoTransaction(rows *sql.Rows, transactions []model.Transaction) ([]model.Transaction,error) {
	for rows.Next() {	
		var t model.Transaction
		err := rows.Scan(&t.ID,&t.From, &t.To, &t.Amount, &t.Time)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions,nil
}

