package controllers

import (
	"database/sql"
	"eqweqr/bdkurach/models"
)

const (
	queryPartOrders = `SELECT id, created_at, model_name, comment, order_status warrantly, client_id from orders, conf_time, summary where client_id=$1 and order_status=$2`
)

// ID          int
// ModelName   string
// Warranty    bool
// Comment     string
// ClientId    int
// OrderStatus string
// ConfTime    string
// Summary     string

func GetOrdersForClientWithStatus(clientId, orderStatus string, db *sql.DB) ([]models.Order, error) {
	var orders []models.Order
	rows, err := db.Query(queryPartOrders, clientId, orderStatus)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		err := rows.Scan(&order.ID, &order.CreatedAt, &order.ModelName, &order.Comment, &order.OrderStatus, &order.Warranty, &order.ClientId, &order.ConfTime, &order.Summary)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
