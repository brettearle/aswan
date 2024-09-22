package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "modernc.org/sqlite"
)

var db *sql.DB

type aswanDB struct {
	path string
	db *sql.DB
}

func (db *aswanDB) createTodo(todo *item) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `INSERT INTO todo (desc, done) VALUES (?,?);`, todo.done, todo.desc,)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not create todo")
	}
	return res, nil
}

func (db *aswanDB) deleteTodo(id int) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `DELETE FROM todo WHERE id = ?;`, id,)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not delete todo")
	}
	return res, nil
}

func (db *aswanDB) updateTodo(todo *item) (sql.Result, error) {
	res, err := db.db.ExecContext(context.Background(), `UPDATE todo SET desc = ?, done = ? WHERE id = ?;`,todo.desc, todo.done, todo.id,)
	if err != nil {
		fmt.Printf("sql Error: %v\n", err)
		return res, errors.New("could not update todo")
	}
	return res, nil
}

func (db *aswanDB) getAllTodos() (*itemList, error) {
	rows, err := db.db.Query("SELECT * FROM todo")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var res itemList
	for rows.Next() {
		var item item
		// TODO rows scan needs error handling
		if err := rows.Scan(&item.id, &item.desc, &item.done); err != nil {
			fmt.Printf("Scan Error: %v\n", err)
		}
		res = append(res, &item)
	}
	return &res, nil
}

func newAswanDB(path string, DB *sql.DB) *aswanDB {
	db := &aswanDB{
		path: path,
		db: DB,
	}
	return db
}

func dbInit(path string) (*aswanDB, error) {
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