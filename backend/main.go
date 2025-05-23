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
	db, err := database.OpenDb("postgres", "postgres://admin:password@localhost:5432/db?sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	serv, err := server.NewServer(secret, "127.0.0.1", 8081, db)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	// служебные
	http.HandleFunc("/api/token/refresh", serv.RefreshToken)
	http.HandleFunc("/api/token/check", serv.CheckToken)
	http.HandleFunc("/api/token/role", serv.GetRole)

	// reg/login
	http.HandleFunc("/api/user/register", serv.Register)
	http.HandleFunc("/api/user/login", serv.Login)
	// http.HandleFunc("/api/order/create", serv.CreateOrderHandler)

	// get data about orders
	// получить все заказы (только для админа)

	// admin
	http.Handle("/api/orders/all", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllOrders), "admin"))
	http.HandleFunc("/api/users/all", serv.GetAllUsersHandler)
	http.HandleFunc("/api/workers/salary", serv.GetWorkersSalaryHandler)
	http.HandleFunc("/api/orders/admin", serv.GetAllOrdersAdminStatusHandler)
	http.HandleFunc("/api/status/change", serv.DeactivateHandler)
	http.HandleFunc("/api/account/new", serv.CreateNewUserHandler)

	// client
	http.HandleFunc("/api/order/create", serv.CreateOrderHandler)
	http.HandleFunc("/api/suggestions", serv.ShowSuggestionByOrderHandler)
	http.HandleFunc("/api/orders", serv.ShowByStatusOrdersHandler)
	http.HandleFunc("/api/order", serv.CancelOrderHandler)
	http.HandleFunc("/api/suggest/pick", serv.AssignWorkerToOrder)
	http.HandleFunc("/api/workers", serv.GetAllWorkersHandler)
	http.HandleFunc("/api/work/types", serv.GetAllTypesHandler)
	http.HandleFunc("/api/order/approve", serv.ApproveWorkHandler)
	http.HandleFunc("/api/order/new", serv.CreateOrderHandler)

	// http.Handle("/api/orders/user", serv.RoleMiddleware(http.HandlerFunc(serv.GetUserOrders), "user"))

	// worker
	http.HandleFunc("/api/worker/orders", serv.GetAllOrdersByStatusWorkerHandler)
	// http.HandleFunc("/api/")
	http.HandleFunc("/api/worker/suggestion", serv.GetAllsugavaitHandler)
	http.HandleFunc("/api/worker/salary/all", serv.AllSalaryHandler)
	// http.HandleFunc("/api/worker/salary", serv.SalaryByTimeHandler)

	http.ListenAndServe(":8080", nil)
}
