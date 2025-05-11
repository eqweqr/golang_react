package main

import (
	"eqweqr/bdkurach/internals/database"
	"eqweqr/bdkurach/internals/password"
	"eqweqr/bdkurach/server"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const (
	DEFAULT_SIZE = 50
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	var err error
	if secret == "" {
		secret, err = password.GenerateSecret(DEFAULT_SIZE)
		if err != nil {
			panic(err)
		}
	}
	db, err := database.OpenDb("postgres", "postgres://mkyong:password@localhost:5432/kurs?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	serv, err := server.NewServer(secret, "127.0.0.1", 8081, db)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	http.HandleFunc("/api/token/refresh", serv.RefreshToken)
	http.HandleFunc("/api/token/check", serv.CheckToken)
	http.HandleFunc("/api/token/role", serv.GetRole)
	http.HandleFunc("/api/user/register", serv.Register)
	http.HandleFunc("/api/user/login", serv.Login)

	http.ListenAndServe(":8081", nil)
}
