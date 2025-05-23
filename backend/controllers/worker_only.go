package controllers

import (
	"database/sql"
	"eqweqr/bdkurach/dto"
	"fmt"
)

func GetAllSuggestionsWorker(worker_id string, db *sql.DB) ([]dto.SugWorker, error) {
	query := `select s.id, p.name, s.summary, s.term, o.name  from suggestions s left join orders o on s.order_id=o.id left join people p on o.client_id=p.id where s.worker_id=$1`
	var sugworker []dto.SugWorker

	rows, err := db.Query(query, worker_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var sug dto.SugWorker
		err := rows.Scan(&sug.ID, &sug.ClientName, &sug.Summary, &sug.Term, &sug.OrderName)
		if err != nil {
			return nil, err
		}
		sugworker = append(sugworker, sug)
	}
	return sugworker, nil

}

func GetAllStatusOrdersWorker(id string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.worker_id=$1 `
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

func CreateSuggestin(order_id, worker_id, summary, term string, db *sql.DB) error {
	query := `insert into suggestions(order_id, worker_id, summary, term) values($1, $2, $3, $4)`
	_, err := db.Exec(query, order_id, worker_id, summary, term)
	if err != nil {
		return fmt.Errorf("Error while creating suggestions: %v", err)
	}
	return nil

}

func GetAllOrdersByTime(worker_id, days string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.worker_id=$1 and o.conf_time < now()-interval '$2 days'`
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, worker_id, days)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func GetTotalSummary(id string, db *sql.DB) (dto.SummaryDto, error) {
	queryDone := `select COALESCE(sum(summary),0) from orders where worker_id=$1 and order_status='done' `

	var sum dto.SummaryDto
	err := db.QueryRow(queryDone, id).Scan(&sum.Summary)
	if err != nil {
		return dto.SummaryDto{}, err
	}

	return sum, nil
}

func GetTotalSummaryByTime(id, days string, db *sql.DB) (dto.SummaryDto, error) {
	queryDone := `select coalesce(sum(summary), 0) from orders where worker_id=$1 and order_status='done' and conf_time < now()-interval $2`

	d := days + " days"
	var sum dto.SummaryDto
	err := db.QueryRow(queryDone, id, d).Scan(&sum.Summary)
	if err != nil {
		return dto.SummaryDto{}, err
	}

	return sum, nil
}

func GetAllOrderStatusWorker(id, status string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.worker_id=$1 and o.order_status=$2 `
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, id, status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

func GetAllSugessiongWorker(id string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name from suggestions s left join orders o on s.order_id=o.id left join people p on o.client_id=p.id left join typework t on o.work_type=t.id left join people p1 on p1.id=s.worker_id  where s.worker_id=$1 `
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}
