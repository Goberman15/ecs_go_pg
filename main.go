package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

type server struct {
	db *sqlx.DB
}

type passenger struct {
	Id        string `db:"passenger_id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Age       int
	FlightNum string `db:"flight_number"`
}

func (s *server) getHandler(w http.ResponseWriter, r *http.Request) {
	var p []passenger

	err := s.db.Select(&p, "SELECT * FROM passenger")
	if err != nil {
		w.WriteHeader(500)
		fmt.Print(err.Error())
		io.WriteString(w, "Query Error")
		return
	}

	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "JSON Stringify Error")
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/json")
}

func main() {
	godotenv.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PARAMS"),
	)

	fmt.Println("dsn =>", dsn)

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		fmt.Println("Fatal: Fail to Connect to DB")
		log.Fatal(err)
	}
	defer db.Close()

	s := &server{
		db,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /test/", s.getHandler)
	log.Fatal(http.ListenAndServe(":8088", mux))
}
