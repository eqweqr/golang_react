package controllers

import "database/sql"

// OrderId  int
// WorkerId int
// Summary  int

const (
	execSuggest = `insert into suggestions(order_id, worker_id, summary, status) values($1, $2, $3, $4) returning id`
)

func CreateNewSuggest(order_id int, worker_id int, summary string, status string, db *sql.DB) (int, error) {
	var id int
	err := db.QueryRow(execSuggest, order_id, worker_id, summary, status).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
