package controllers

import (
	"database/sql"
	"eqweqr/bdkurach/models"
)

const (
	execInsert = `INSERT INTO people(name, role, email, phone, password, status) values ($1, $2, $3, $4, $5, 'active') returning id`
	queryPass  = `SELECT id, name, email, phone, role, password from people where email=$1`
	checkPhone = `SELECT id from people where phone=$1`
)

func CheckPhoneExists(phone string, db *sql.DB) (bool, error) {
	rows, err := db.Query(checkPhone, phone)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func AddNewUser(username string, phone string, email string, pass string, secret string, db *sql.DB) (int, error) {

	var id int
	// encrypted_pass, err := password.HashPassword(pass)
	// if err != nil {
	// 	return 0, err
	// }

	err := db.QueryRow(execInsert, username, "client", email, phone, pass).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func GetUserByPhone(phone string, db *sql.DB) (*models.User, error) {
	var user models.User
	err := db.QueryRow(queryPass, phone).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Role, &user.Pass)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
