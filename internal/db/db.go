package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type AswanDB struct {
	Path string
	Instance   *sql.DB
}

func (db *AswanDB) CreateTodo(desc string, done bool) (sql.Result, error) {
	res, err := db.Instance.ExecContext(context.Background(), `INSERT INTO todo (desc, done) VALUES (?,?);`, desc, done)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not create todo")
	}
	return res, nil
}

func (db *AswanDB) DeleteTodo(id int) (sql.Result, error) {
	res, err := db.Instance.ExecContext(context.Background(), `DELETE FROM todo WHERE id = ?;`, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not delete todo")
	}
	return res, nil
}

func (db *AswanDB) UpdateTodo(id int, desc string, done bool) (sql.Result, error) {
	res, err := db.Instance.ExecContext(context.Background(), `UPDATE todo SET desc = ?, done = ? WHERE id = ?;`, desc, done, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not update todo")
	}
	return res, nil
}

func (db *AswanDB) GetAllTodos() (*sql.Rows, error) {
	rows, err := db.Instance.Query("SELECT * FROM todo")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return rows, nil
}

func NewAswanDB(path string, DB *sql.DB) *AswanDB {
	db := &AswanDB{
		Path: path,
		Instance:   DB,
	}
	return db
}

func GetDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// Create the .aswan directory path
	dbDir := filepath.Join(homeDir, ".aswan")

	// Ensure the .aswan directory exists
	if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
		panic(err)
	}

	// Define the SQLite database file path
	dbFile := "aswan.db"
	return filepath.Join(dbDir, dbFile)
}

func DbInit(path string) (*AswanDB, error) {
	//Init Sqlite
	db, err := sql.Open("sqlite", path)
	if err != nil {
		fmt.Printf("Failed to open db %v\n", err)
		return nil, err
	}
	//Init Tables
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
		return nil, err
	}
	//Return new aswan db
	r := NewAswanDB(path, db)
	return r, nil
}
