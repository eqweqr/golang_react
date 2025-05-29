package controllers

import (
	"database/sql"
	"eqweqr/bdkurach/dto"
	"fmt"
)

const (
	getAllPageOrders  = `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email, o.term from orders o join people p on o.client_id=p.id join people on p1 o.worker_id=p1.id limit $1 offset $1*$2`
	getAllClientsName = `select id, name from people where role='client' `
	getAllClientsFull = `select id, name, email, phone from people where role='client'`
	limit             = 20
)

// get all orders
func GetOrdersAdminByPage(page string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	var orders []dto.FullOrderDTO
	rows, err := db.Query(getAllPageOrders, limit, page)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.Term)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// get all clients
func GetAllClientFull(db *sql.DB) ([]dto.FullUserDTO, error) {
	var clients []dto.FullUserDTO
	rows, err := db.Query(getAllClientsFull)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var client dto.FullUserDTO
		err := rows.Scan(&client.Id, &client.Name, &client.Email, &client.Phone)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func CheckIsActive(id string, db *sql.DB) error {
	query := `select 1 from people where id=$1 and status='active'`
	var v int
	err := db.QueryRow(query, id).Scan(&v)
	if err != nil {
		return err
	}
	return nil
}

// get all workers names
func GetAllClientName(db *sql.DB) ([]dto.NameUserDTO, error) {
	var clients []dto.NameUserDTO
	rows, err := db.Query(getAllClientsName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var client dto.NameUserDTO
		err := rows.Scan(&client.Id, &client.Name)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

// create new account
func CreateNewAccount(username string, phone string, email string, pass string, role string, db *sql.DB) error {
	query := `INSERT INTO people(name, role, email, phone, password) values ($1, $2, $3, $4, $5) returning id`

	// encrypted_pass, err := password.HashPassword(pass)
	// if err != nil {
	// 	return 0, err
	// }
	var tmp interface{}

	err := db.QueryRow(query, username, role, email, phone, pass).Scan(&tmp)

	if err != nil {
		return fmt.Errorf("Error occured while adding new user: %v", err)
	}

	return nil
}

// create new account with return status
func CreateNewAccountWithStatus(username string, phone string, email string, pass string, role string, db *sql.DB) (dto.StatusDTO, error) {
	query := `INSERT INTO people(name, role, email, phone, password, status) values ($1, $2, $3, $4, $5, 'active') returning id`

	// encrypted_pass, err := password.HashPassword(pass)
	// if err != nil {
	// 	return 0, err
	// }
	var tmp dto.StatusDTO

	err := db.QueryRow(query, username, role, email, phone, pass).Scan(&tmp.Id)
	tmp.Status = "active"

	if err != nil {
		return dto.StatusDTO{}, fmt.Errorf("Error occured while adding new user: %v", err)
	}

	return tmp, nil
}

// могут не отображаться заказы на которых ен был
func GetOrders(days string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name, o.term from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.conf_time < now()-interval '$1 days'`
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, days)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType, &order.Term)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

// status = pending/ processing/ done
func GetAllStatusOrdersAdmin(status string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name, o.term from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.order_status=$1`
	// db.Query(query, id, status)
	var orders []dto.FullOrderDTO
	rows, err := db.Query(query, status)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.FullOrderDTO
		err := rows.Scan(&order.OrderId, &order.CreatedAt, &order.ModelName, &order.Warranty, &order.Comment, &order.OrderStatus, &order.ConfTime, &order.Summary, &order.ClientName, &order.ClientPhone, &order.ClientEmail, &order.WorkerName, &order.WorkerPhone, &order.WorkerEmail, &order.WorkType, &order.Term)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil

}

// get worker и выплаты
func GetWorkersSalary(db *sql.DB) ([]dto.UserSalary, error) {
	query := `select p.name, p.phone, sum(o.summary), p.email from people p left join orders o on (p.id=o.worker_id and o.order_status='done') group by p.id having p.role='worker'`
	var tmps []dto.UserSalary
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tmp dto.UserSalary
		err := rows.Scan(&tmp.Name, &tmp.Phone, &tmp.Summary, &tmp.Email)
		if err != nil {
			return nil, err
		}
		tmps = append(tmps, tmp)
	}
	return tmps, nil

}
func DeactivateHandler(id string, status string, db *sql.DB) error {
	query := `update people set status=$1 where id=$2`

	_, err := db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("while deactivating %v", err)
	}
	return nil
}

// all users
func GetAllUsers(db *sql.DB) ([]dto.UserRoleDTO, error) {
	query := `select id, name, phone, email, role, status from people`
	var tmps []dto.UserRoleDTO
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tmp dto.UserRoleDTO
		err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Phone, &tmp.Email, &tmp.Role, &tmp.Status)
		if err != nil {
			return nil, err
		}
		tmps = append(tmps, tmp)
	}
	return tmps, nil
}

// disable role
