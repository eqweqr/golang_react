package controllers

import (
	"database/sql"
	"eqweqr/bdkurach/dto"
	"fmt"
	"log"
)

const (
	getAllWorkersName          = `select id, name from people where role='worker' `
	getAllWorkersFull          = `select id, name, email, phone from people where role='worker'`
	deleteOrderRequest         = `delete from orders where id=$1 returning id`
	cascadeDeleteAllSuggestion = `delete from suggestions where order_id=$1`
	cancelSuggestions          = `delete from suggestions where order_id=$1`
)

// создать заказ
func CreateNewOrder(model_name, comment, client_id, order_status string, warranty bool, work_type string, worker string, db *sql.DB) (dto.Ids, error) {
	workerId := `select id from typework where name=$1`
	deviceId := `select id from people where name=$1`
	createNewOrder := `INSERT INTO orders(model_name, comment, client_id, order_status, warranty, work_type, worker_id) values ($1, $2, $3, $4, $5, $6, $7) returning id`
	createNewOrder1 := `INSERT INTO orders(model_name, comment, client_id, order_status, warranty, work_type) values ($1, $2, $3, $4, $5, $6) returning id`

	var tt dto.Ids
	ext := true
	if worker == "Любой" {
		ext = false
	}
	var dId int
	log.Println(work_type)
	err := db.QueryRow(deviceId, worker).Scan(&dId)
	if err != nil {
		return dto.Ids{}, fmt.Errorf("finding device id: %v", err)
	}

	if ext {
		var wId int
		err = db.QueryRow(workerId, work_type).Scan(&wId)
		if err != nil {
			return dto.Ids{}, fmt.Errorf("finding worker id: %v", err)
		}
		err = db.QueryRow(createNewOrder, model_name, comment, client_id, order_status, warranty, dId, wId).Scan(&tt.Id)

		if err != nil {
			return dto.Ids{}, fmt.Errorf("creating new order: %v", err)
		}

		err = db.QueryRow(`select phone from people where name=$1`, worker).Scan(&tt.Phone)

		if err != nil {
			return dto.Ids{}, fmt.Errorf("creating new order: %v", err)
		}
		return tt, nil
	}
	err = db.QueryRow(createNewOrder1, model_name, comment, client_id, order_status, warranty, dId).Scan(&tt.Id)

	if err != nil {
		return dto.Ids{}, fmt.Errorf("creating new order: %v", err)
	}

	err = db.QueryRow(`select phone from people where name=$1`, worker).Scan(&tt.Phone)

	if err != nil {
		return dto.Ids{}, fmt.Errorf("creating new order: %v", err)
	}

	return tt, nil
}

// get all workers full
func GetAllWorkerFull(db *sql.DB) ([]dto.FullUserDTO, error) {
	var clients []dto.FullUserDTO
	rows, err := db.Query(getAllWorkersFull)
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

// get all workers names
func GetAllWorkersName(db *sql.DB) ([]dto.NameUserDTO, error) {
	var clients []dto.NameUserDTO
	rows, err := db.Query(getAllWorkersName)
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

// check if status is pending else user cant cancel order
func CheckOrderStatus(id, status string, tx *sql.Tx) error {
	query := `select 1 from orders where id=$1 and order_status=$2`
	row := tx.QueryRow(query, id, status)
	var ans interface{}
	err := row.Scan(&ans)

	if err != nil {
		return err
	}
	return nil
}

// delete order(only if it not approved)
func CancelOrder(id string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error on startign tx: t%v", err)
	}
	defer tx.Rollback()

	err = CheckOrderStatus(id, "pending", tx)
	if err != nil {
		return fmt.Errorf("error on checking status: t%v", err)
	}

	_, err = tx.Exec(cascadeDeleteAllSuggestion, id)
	if err != nil {
		return fmt.Errorf("error on cascade deleting suggestions: t%v", err)
	}

	row := tx.QueryRow(deleteOrderRequest, id)

	var rid interface{}
	err = row.Scan(&rid)

	if err != nil {
		return fmt.Errorf("error on cascade deleting order: t%v", err)
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("error on commit: t%v", err)
	}

	return nil
}

// отменить одно предложение может выполнить как клиент так и работник
func DeleteSuggestion(id int, db *sql.DB) error {
	_, err := db.Exec(cancelSuggestions, id)
	return err
}

// позвоялет обновлять статус заказа
func UpdateOrderStatus(id int, status string, db *sql.DB) error {
	updateQuery := `update orders set status=$2 where id=$1`
	_, err := db.Exec(updateQuery, id, status)
	return err
}

// check order belong to user
func CheckOrderBelong(order_id, id string, db *sql.DB) error {
	query := `select 1 from orders where id=$1 and client_id = $2`
	var v interface{}
	err := db.QueryRow(query, order_id, id).Scan(&v)
	if err != nil {
		return fmt.Errorf("Error occured while check belong order to user: %v", err)
	}
	return nil
}

// get all suggestions for user
func GetAllSuggestions(user_id string, db *sql.DB) ([]dto.ExtendedSuggestionDTO, error) {
	query := `select s.id, s.order_id, p.name, s.summary, s.term, o.model_name, p.phone from suggestions s left join people p on s.worker_id=p.id left join orders o on s.order_id=o.id where o.client_id=$1`
	var suggestions []dto.ExtendedSuggestionDTO
	rows, err := db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var suggestion dto.ExtendedSuggestionDTO
		err := rows.Scan(&suggestion.ID, &suggestion.OID, &suggestion.WorkerName, &suggestion.Summary, &suggestion.Term, &suggestion.WorkName, &suggestion.Phone)
		if err != nil {
			return nil, err
		}
		suggestions = append(suggestions, suggestion)
	}
	return suggestions, nil

}

// удалить все предложения связанные с заказом
func DeleteAllRelatedSuggestion(id int, tx *sql.Tx) error {
	_, err := tx.Exec(cancelSuggestions, id)
	return err
}

type partSug struct {
	Summary  int
	OrderId  int
	WorkerId int
	Term     int
}

// назначить работнику заказ
func AssignWorkerToOrder(sugId string, db *sql.DB) error {
	execUpdate := `update orders set worker_id=$1, summary=$2, term=$3, conf_time=now(), order_status='processing' where id=$4`
	getValFromSug := `select id, order_id, worker_id, summary from suggestions where id=$1`

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var suggest partSug
	err = tx.QueryRow(getValFromSug, sugId).Scan(&suggest.OrderId, &suggest.WorkerId, &suggest.Summary, &suggest.Term)
	if err != nil {
		return fmt.Errorf("Error while getting suggestion on id : %v", err)
	}

	_, err = tx.Exec(execUpdate, suggest.WorkerId, suggest.Summary, suggest.Term, suggest.OrderId)
	if err != nil {
		return fmt.Errorf("Error while changing order : %v", err)
	}

	err = DeleteAllRelatedSuggestion(suggest.OrderId, tx)
	if err != nil {
		return fmt.Errorf("Error while deleting all related suggestions : %v", err)
	}

	return tx.Commit()
}

// check if suggestion belong to user`s order
func SuggestionOrder(sug_id, user_id string, db *sql.DB) error {
	query := `select 1 from suggestions s join orders o on s.order_id=o.id where s.id=$1 and o.client_id=$2`
	var tmp interface{}
	err := db.QueryRow(query, sug_id, user_id).Scan(&tmp)
	if err != nil {
		return fmt.Errorf("err on suggestion doesnt belong user`s order: %v", err)
	}
	return nil
}

// status = pending/ processing/ done
func GetAllStatusOrders(id, status string, db *sql.DB) ([]dto.FullOrderDTO, error) {
	query := `select o.id, o.created_at, o.model_name, o.warranty, o.comment, o.order_status, o.conf_time, o.summary, p.name, p.phone, p.email, p1.name, p1.phone, p1.email , t.name from orders o left join people p on o.client_id=p.id left join people p1 on o.worker_id=p1.id left join typework t on o.work_type=t.id where o.client_id=$1 and o.order_status=$2 `
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

func GetAllTypes(db *sql.DB) ([]dto.WorkTypeDTO, error) {
	query := `select id, name from typework`

	var orders []dto.WorkTypeDTO
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order dto.WorkTypeDTO
		err := rows.Scan(&order.Id, &order.Name)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func ApproveWork(order_id string, db *sql.DB) error {
	query := `update orders set order_status='done' where id=$1 and order_status='processing' returning id`
	var tmp interface{}
	err := db.QueryRow(query, order_id).Scan(&tmp)
	if err != nil {
		return fmt.Errorf("Error occured while updating work status: %v", err)
	}
	return nil
}
