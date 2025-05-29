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
	"strconv"
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

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	type t struct {
		Name  string `json:"username"`
		Pass  string `json:"password"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	var tmp t
	err := d.Decode(&tmp)
	log.Println(&tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if len(tmp.Name) < 6 || len(tmp.Phone) != 11 || len(tmp.Pass) < 6 || !strings.Contains(tmp.Email, "@") {
	// log.Println("invalid input format")
	// w.WriteHeader(http.StatusUnauthorized)
	// return
	// }

	// ok, err := controllers.CheckPhoneExists(tmp.Phone, server.DB)
	// if ok || err != nil {
	// 	log.Println("phone number already exists")
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

	encrypted_pass, err := password.HashPassword(tmp.Pass)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id, err := controllers.AddNewUser(tmp.Name, tmp.Phone, tmp.Email, encrypted_pass, server.Secret, server.DB)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := jwttoken.CreateToken(server.Secret, id, tmp.Name, standart_role)
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

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	type t struct {
		Name  string `json:"username"`
		Pass  string `json:"password"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}

	var tmp t
	err := d.Decode(&tmp)
	log.Println(&tmp)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	phone := tmp.Email
	pass := tmp.Pass

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

	if err = controllers.CheckIsActive(strconv.Itoa(user.ID), server.DB); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	regdto := dto.NewRegisterDTO(token)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(regdto)
}
