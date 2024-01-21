package main

import (
	"ewallet/internal/config"
	"ewallet/pkg/app"
	postgres "ewallet/pkg/storage/postgres"

	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)


func main(){
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	DB_URL := viper.GetString("db.url")

	store, err := postgres.NewPostgres(DB_URL)
	if err != nil {
		log.Fatal("err: failed connection to postgres")
	}

	defer store.DB.Close()

	if err = store.InitTables(); err != nil {
		log.Fatal("err: failed to init tables",err)
	}

	listenAddr := viper.GetString("listenAddr")
	server := app.NewAPIServer(listenAddr, store)
	server.RunServer()

}
