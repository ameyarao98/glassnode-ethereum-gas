package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx"
)

var conn *pgx.ConnPool // global postgres connection pool
var dumpDate time.Time = time.Date(2020, 9, 7, 0, 0, 0, 0, time.UTC)

type eoaFee struct {
	Timestamp int64   `json:"t"`
	Fees      float64 `json:"v"`
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
	var err error
	conn, err = pgx.NewConnPool(pgxConnPoolConfig)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Connected to postgres")
	}
	http.HandleFunc("/eoa-fees/", GetEOAFees)
	http.ListenAndServe(":8080", nil)
}

func ContractAddresses() ([]string, error) {
	contracts := make([]string, 0)
	rows, err := conn.Query("SELECT \"address\" FROM \"contracts\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var contract string
		if err := rows.Scan(&contract); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}
	return contracts, nil
}

func HourlyGasData() ([]eoaFee, error) {
	nonEOAAddresses, err := ContractAddresses()
	if err != nil {
		return nil, err
	}
	nonEOAAddresses = append(nonEOAAddresses, "0x0000000000000000000000000000000000000000")
	rows, err := conn.Query("SELECT EXTRACT(HOUR FROM \"block_time\") AS \"hour\", SUM((\"gas_price\"/1000000000000000000)*\"gas_used\") FROM \"transactions\" WHERE \"to\" <> ANY ($1) AND \"from\" <> ANY ($1) GROUP BY \"hour\" ORDER BY \"hour\"", nonEOAAddresses)
	if err != nil {
		return nil, err
	}
	fees := make([]eoaFee, 0)
	defer rows.Close()
	for rows.Next() {
		var h int
		var v float64
		if err := rows.Scan(&h, &v); err != nil {
			return nil, err
		}
		fees = append(fees, eoaFee{hourToUnixTimeStamp(h), math.Round(v*100) / 100}) // rounding in SQL causes issue when scanned into golang's float
	}
	return fees, nil
}

func hourToUnixTimeStamp(hour int) int64 {
	return time.Date(dumpDate.Year(), dumpDate.Month(), dumpDate.Day(), hour, dumpDate.Minute(), dumpDate.Second(), dumpDate.Nanosecond(), dumpDate.UTC().Location()).Unix()
}

func GetEOAFees(w http.ResponseWriter, r *http.Request) {
	fees, err := HourlyGasData()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(fees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
