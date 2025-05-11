package database

import "database/sql"

func OpenDb(schema string, url string) (*sql.DB, error) {
	db, err := sql.Open(schema, url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
