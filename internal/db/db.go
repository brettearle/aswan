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
	Path     string
	Instance *sql.DB
}

func (db *AswanDB) CreateTodo(desc string, done bool, doneTime string, board string) (sql.Result, error) {
	if len(desc) == 0 {
		panic("Description is empty")
	}
	if len(doneTime) == 0 {
		panic("DoneTime is empty")
	}
	if len(board) == 0 {
		panic("Board is empty")
	}

	result, err := db.Instance.ExecContext(context.Background(), `INSERT INTO todo (desc, done, doneTime, board) VALUES (?,?,?,?);`, desc, done, doneTime, board)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return result, errors.New("could not create todo")
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected > 1 {
		panic("Create affects more than 1 row")
	}

	newID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Error: create todo Last inserted ID \n%v", err)
	}
	if newID < 0 {
		panic("ID of created todo cannot be found")
	}
	return result, nil
}

func (db *AswanDB) DeleteTodo(id int) (sql.Result, error) {
	if id < 0 {
		panic("Delete todo requires id above zero")
	}
	result, err := db.Instance.ExecContext(context.Background(), `DELETE FROM todo WHERE id = ?;`, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return result, errors.New("could not delete todo")
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected > 1 {
		panic("Delete affects more than 1 row")
	}
	return result, nil
}

func (db *AswanDB) UpdateTodo(id int, desc string, done bool, doneTime string, board string) (sql.Result, error) {
	if id < 0 {
		panic("ID must be 0 or more")
	}
	if len(desc) == 0 {
		panic("Description is empty")
	}
	if len(doneTime) == 0 {
		panic("DoneTime is empty")
	}
	if len(board) == 0 {
		panic("Board is empty")
	}

	result, err := db.Instance.ExecContext(context.Background(), `UPDATE todo SET desc = ?, done = ?, doneTime = ?, board = ? WHERE id = ?;`, desc, done, doneTime, board, id)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return result, errors.New("could not update todo")
	}

	rowsAffected, err := result.RowsAffected()
	if rowsAffected > 1 {
		panic("Update affects more than 1 row")
	}
	return result, nil
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
	if len(path) < 1 {
		panic("Path handed to NewAswanDB is empty")
	}
	err := DB.Ping()
	if err != nil {
		panic("Unable to establish connection on DB handed to NewAswanDB")
	}

	db := &AswanDB{
		Path:     path,
		Instance: DB,
	}

	if path != db.Path {
		panic("Path discrepency in NewAswanDB")
	}
	err = db.Instance.Ping()
	if err != nil {
		fmt.Printf("ERR: %v", err)
		panic("Can't ping instance in NewAswanDB")
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
	result := filepath.Join(dbDir, dbFile)
	if len(result) < 1 {
		panic("Length of filepath to DB not initialised correctly")
	}
	return result
}

func DbInit(path string) (*AswanDB, error) {
	if len(path) < 1 {
		panic("Length of filepath to DB not initialised correctly")
	}
	//Init Sqlite
	db, err := sql.Open("sqlite", path)
	if err != nil {
		fmt.Printf("Failed to open db %v\n", err)
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		panic("Could not connect to DB before initialising tables")
	}
	//Init Tables
	_, err = db.ExecContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS todo (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			desc TEXT NOT NULL, 
			done BOOLEAN NOT NULL,
			doneTime TEXT NOT NULL,
			board TEXT NOT NULL
		)`,
	)
	if err != nil {
		fmt.Printf("Failed to create: %v\n", err)
		return nil, err
	}
	//Return new aswan db
	result := NewAswanDB(path, db)
	if len(result.Path) < 1 {
		panic("path not set on NewAswanDB")
	}
	err = result.Instance.Ping()
	if err != nil {
		panic("unable to establish connection on NewAswanDB")
	}
	return result, nil
}
