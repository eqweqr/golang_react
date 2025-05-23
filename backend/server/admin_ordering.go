package server

import (
	"encoding/json"
	"eqweqr/bdkurach/controllers"
	"eqweqr/bdkurach/internals/password"
	"log"
	"net/http"
	"strings"
)

func (server *Server) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	ordersInfo, err := controllers.GetOrdersAdminByPage(page, server.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(ordersInfo)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (server *Server) GetAllOrdersAdminStatusHandler(w http.ResponseWriter, r *http.Request) {
	// id := r.URL.Query().Get("id")
	status := r.URL.Query().Get("status")

	val, err := controllers.GetAllStatusOrdersAdmin(status, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(val)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

func (server *Server) GetWorkersSalaryHandler(w http.ResponseWriter, r *http.Request) {
	slrys, err := controllers.GetWorkersSalary(server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(slrys)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func (server *Server) DeactivateHandler(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	id := r.URL.Query().Get("id")
	err := controllers.DeactivateHandler(id, status, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// получить данные по всем пользователям
func (server *Server) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := controllers.GetAllUsers(server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// создать нового пользователя любой роли
func (server *Server) CreateNewUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	type regnew struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
		Role  string `json:"role"`
		Pass  string `json:"pass"`
	}

	var tmp regnew
	err := d.Decode(&tmp)
	log.Println(&tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	username := tmp.Name
	pass := tmp.Pass
	phone := tmp.Phone
	email := tmp.Email
	role := tmp.Role

	if len(username) < 6 || len(phone) != 11 || len(pass) < 6 || !strings.Contains(email, "@") {
		log.Println("invalid input format")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ok, err := controllers.CheckPhoneExists(phone, server.DB)
	if ok || err != nil {
		log.Println("phone number already exists")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	encrypted_pass, err := password.HashPassword(pass)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	val, err := controllers.CreateNewAccountWithStatus(username, phone, email, encrypted_pass, role, server.DB)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = json.NewEncoder(w).Encode(val); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

}

// получить заказы по времени и пользователю
// func (server *Server)

// посмотреть сколько должны получить работники
