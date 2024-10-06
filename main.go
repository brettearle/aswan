package main

import (
	"fmt"
	"os"

	"github.com/brettearle/aswan/internal/db"
	"github.com/brettearle/aswan/internal/flagger"
	"github.com/brettearle/aswan/internal/todo"
)

func run(db *db.AswanDB) (bool, error) {
	// Initial State
	todosList, err := todo.NewTodoList().Populate(db)
	if err != nil {
		fmt.Println("\nfailed to get todos")
		return false, err
	}
	//Flag Decleration
	flags, err := flagger.FlagService(os.Args, db, todosList)
	if err != nil {
		fmt.Println("\nFailed to init flags")
		return false, err
	}
	//Handlers
	if *flags.Tick {
		todosList, err = flagger.TickHandler(todosList, flags, db)
		if err != nil {
			return false, err
		}
	}

	if *flags.New {
		todosList, err = flagger.NewHandler(todosList, flags, db)
		if err != nil {
			return false, err
		}
	}

	if *flags.Delete {
		todosList, err = flagger.DeleteHandler(todosList, flags, db)
		if err != nil {
			return false, err
		}
	}

	if *flags.Clear {
		_, err = flagger.ClearHandler(todosList, flags, db)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func main() {
	//DB Initialization
	DB, err := db.DbInit(db.GetDBPath())
	if err != nil {
		panic("no DB able to be initialized")
	}
	defer DB.Instance.Close()

	//RUN RUN RUN
	_, err = run(DB)
	if err != nil {
		fmt.Printf("\nRun failed with: %v", err)
	}
}
