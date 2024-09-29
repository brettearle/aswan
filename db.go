package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type aswanDB struct {
	path string
	db   *sql.DB
}

func (db *aswanDB) createTodo(desc string, done bool) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `INSERT INTO todo (desc, done) VALUES (?,?);`, desc, done)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not create todo")
	}
	return res, nil
}

func (db *aswanDB) deleteTodo(id int) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `DELETE FROM todo WHERE id = ?;`, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not delete todo")
	}
	return res, nil
}

func (db *aswanDB) updateTodo(id int, desc string, done bool) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `UPDATE todo SET desc = ?, done = ? WHERE id = ?;`, desc, done, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not update todo")
	}
	return res, nil
}

func (db *aswanDB) getAllTodos() (*sql.Rows, error) {
	rows, err := db.db.Query("SELECT * FROM todo")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	return rows, nil
}

func newAswanDB(path string, DB *sql.DB) *aswanDB {
	db := &aswanDB{
		path: path,
		db:   DB,
	}
	return db
}

func getDBPath() string {
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

func dbInit(path string) (*aswanDB, error) {
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
	r := newAswanDB(path, db)
	return r, nil
}
