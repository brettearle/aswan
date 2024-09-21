package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type AswanDB struct {
	Path string
	DB *sql.DB
}

func newAswanDB(path string, DB *sql.DB) *AswanDB {
	db := &AswanDB{
		Path: path,
		DB: DB,
	}
	return db
}

func Init(path string) (*AswanDB, error) {
	var err error
	db, err = sql.Open("sqlite", path)
	if err != nil {
		fmt.Printf("Failed to open db %v\n", err)
		return nil,err
	}
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS todo (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			desc TEXT NOT NULL, 
			done BOOLEAN NOT NULL 
		)`,
	)
	if err != nil {
		fmt.Printf("Failed to create %v\n", err)
		return nil,err
	}
	r := newAswanDB(path, db)	
	return r, nil
}