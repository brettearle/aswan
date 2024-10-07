package main

import (
	"bufio"
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
	CWD, err := os.ReadDir("./")
	if err != nil {
		fmt.Println("Could no read current Dir")
		os.Exit(1)
	}

	dbExists := false
	for _, file := range CWD {
		if file.Name() == ".aswan" {
			dbExists = true
		}
	}

	if !dbExists {
		//DB Initialization
		reader := bufio.NewReader(os.Stdin)
		fmt.Println("Would you like to start a list in this directory? y/n")
		key, err := reader.ReadString('\n')
		if err != nil || key != "y\n" {
			fmt.Println("DB not initialised")
			os.Exit(1)
		}
	}

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
