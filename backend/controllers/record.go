package controllers

import (
	"context"
	"database/sql"
	"eqweqr/bdkurach/models"
	"fmt"
)

func rangeQuery(db sql.DB, query string, ctx context.Context) ([]models.Record, error) {
	var records []models.Record
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		fmt.Println("print another record")
		var record models.Record
		err := rows.Scan(&record.ID, &record.Name, &record.Age)
		if err != nil {
			return nil, err
		}
		fmt.Print(record)
		records = append(records, record)
	}

	return records, nil
}
