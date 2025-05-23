package server

import (
	"encoding/json"
	"eqweqr/bdkurach/controllers"
	"eqweqr/bdkurach/dto"
	"eqweqr/bdkurach/internals/jwttoken"
	"eqweqr/bdkurach/internals/password"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	standart_role = "client"
)

func (server *Server) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	token, err := jwttoken.CreateToken(server.Secret, 100, "root", "user")
	fmt.Println(token)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(token))
}

func (server *Server) CheckToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	bearer := r.Header.Get("Authorization")
	jwtToken := jwttoken.GetToken(bearer)
	log.Println(jwtToken)
	_, err := jwttoken.ParseToken(jwtToken, server.Secret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (server *Server) GetRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	bearer := r.Header.Get("Authorization")
	jwtToken := jwttoken.GetToken(bearer)
	log.Println(jwtToken)
	token, err := jwttoken.ParseToken(jwtToken, server.Secret)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	role := jwttoken.GetRoles(token)
	roleDTO := dto.NewRoleDTO(role)
	json.NewEncoder(w).Encode(roleDTO)
}

func (server *Server) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	pass := r.FormValue("password")
	phone := r.FormValue("phone")
	email := r.FormValue("email")

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

	id, err := controllers.AddNewUser(username, phone, email, encrypted_pass, server.Secret, server.DB)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := jwttoken.CreateToken(server.Secret, id, username, standart_role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	regdto := dto.NewRegisterDTO(token)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(regdto)
}

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	phone := r.FormValue("phone")
	pass := r.FormValue("password")

	user, err := controllers.GetUserByPhone(phone, server.DB)

	if err != nil {
		log.Printf("error while processing query from db with phone number: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !password.VerifyPassword(pass, user.Pass) {
		log.Println("Passwords not equal")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := jwttoken.CreateToken(server.Secret, user.ID, user.Name, user.Role)

	if err != nil {
		log.Println("error while creating token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	regdto := dto.NewRegisterDTO(token)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(regdto)
}
