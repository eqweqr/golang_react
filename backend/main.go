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
	http.Handle("/api/users/all", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllUsersHandler), "admin"))
	http.Handle("/api/workers/salary", serv.RoleMiddleware(http.HandlerFunc(serv.GetWorkersSalaryHandler), "admin"))
	http.Handle("/api/orders/admin", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllOrdersAdminStatusHandler), "admin"))
	http.Handle("/api/status/change", serv.RoleMiddleware(http.HandlerFunc(serv.DeactivateHandler), "admin"))
	http.Handle("/api/account/new", serv.RoleMiddleware(http.HandlerFunc(serv.CreateNewUserHandler), "admin"))

	// client
	http.Handle("/api/order/create", serv.RoleMiddleware(http.HandlerFunc(serv.CreateOrderHandler), "client"))
	http.Handle("/api/suggestions", serv.RoleMiddleware(http.HandlerFunc(serv.ShowSuggestionByOrderHandler), "client"))
	http.Handle("/api/worker/agree", serv.RoleMiddleware(http.HandlerFunc(serv.AgreeOrderHandler), "client"))
	http.Handle("/api/client/waiting", serv.RoleMiddleware(http.HandlerFunc(serv.ShowByStatusOrdersIndiHandler), "client"))
	http.Handle("/api/orders", serv.RoleMiddleware(http.HandlerFunc(serv.ShowByStatusOrdersHandler), "client"))
	http.Handle("/api/order", serv.RoleMiddleware(http.HandlerFunc(serv.CancelOrderHandler), "client"))
	http.Handle("/api/suggest/pick", serv.RoleMiddleware(http.HandlerFunc(serv.AssignWorkerToOrder), "client"))
	http.Handle("/api/workers", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllWorkersHandler), "client"))
	http.Handle("/api/work/types", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllTypesHandler), "client"))
	http.Handle("/api/order/approve", serv.RoleMiddleware(http.HandlerFunc(serv.ApproveWorkHandler), "client"))
	http.Handle("/api/order/new", serv.RoleMiddleware(http.HandlerFunc(serv.CreateOrderHandler), "client"))

	// http.Handle("/api/orders/user", serv.RoleMiddleware(http.HandlerFunc(serv.GetUserOrders), "user"))

	// worker
	http.Handle("/api/worker/orders", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllOrdersByStatusWorkerHandler), "worker"))
	http.Handle("/api/worker/suggestion", serv.RoleMiddleware(http.HandlerFunc(serv.GetAllsugavaitHandler), "worker"))
	http.Handle("/api/worker/salary/all", serv.RoleMiddleware(http.HandlerFunc(serv.AllSalaryHandler), "worker"))
	http.Handle("/api/worker/salary", serv.RoleMiddleware(http.HandlerFunc(serv.SalaryByTimeHandler), "worker"))
	http.Handle("/api/worker/my/suggestion", serv.RoleMiddleware(http.HandlerFunc(serv.HandleMySuggestionsHandler), "worker"))
	http.Handle("/api/worker/all/suggestion", serv.RoleMiddleware(http.HandlerFunc(serv.HandleAllWorkerSuggestionHandler), "worker"))
	http.Handle("/api/worker/all/suggestion/send", serv.RoleMiddleware(http.HandlerFunc(serv.CreateNewSuggestHandler), "worker"))

	http.ListenAndServe(":8080", nil)
}
