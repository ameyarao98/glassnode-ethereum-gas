package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx"
)

type eoaFees struct {
	Timestamp int64 `json:"t"`
	Fees      int32 `json:"v"`
}

func main() {
	pgxConfig := pgx.ConnConfig{
		Host:     "database",
		Port:     5432,
		Database: os.Getenv("POSTGRES_DB"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}
	pgxConnPoolConfig := pgx.ConnPoolConfig{ConnConfig: pgxConfig}
	_, err := pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to postgres")
	}
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":8080", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Yo")
}
