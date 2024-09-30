package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/brettearle/aswan/internal/db"
	"github.com/brettearle/aswan/internal/todo"
)

type todoFlags struct {
	nameArg    string
	itemFlags  *flag.FlagSet
	newFlag    *bool
	tickFlag   *bool
	deleteFlag *bool
	clearFlag  *bool
}

func newTodoFlags(desc string) *todoFlags {
	iF := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	return &todoFlags{
		itemFlags:  iF,
		newFlag:    iF.Bool("n", false, "New todo"),
		tickFlag:   iF.Bool("t", false, "Completes a todo"),
		deleteFlag: iF.Bool("d", false, "Deletes a todo"),
		clearFlag:  iF.Bool("clear", false, "Deletes all todos"),
		nameArg:    desc,
	}
}

func flagService(
	args []string,
	db *db.AswanDB,
	todosList *todo.Todolist,
) (*todoFlags, error) {

	var err error
	//Commands
	commands := args
	if len(commands) == 1 {
		todo.RenderTodos(db, callClear)
		return newTodoFlags(""), err
	}

	flags := newTodoFlags(args[1])

	_, clear, flagFirst := strings.Cut(commands[1], "-")

	//Handles structure
	if flagFirst && clear != "clear" {
		fmt.Println("\nPlease provide item string first")
		fmt.Println("example: `CMD> aswan 'im a item' -n -t`")
		flags.itemFlags.PrintDefaults()
		return nil, err
	}

	if len(commands) > 1 {
		switch commands[1] {
		//cases for commands go here
		case "help":
			fmt.Println("\nhelp not implemented")
			return flags, nil
		case "dbPath":
			fmt.Printf("\n%s\n", db.Path)
			return flags, nil
		case "timer":
			fmt.Println("\ntimer not yet implemented")
			return flags, nil
		case "rmDone":
			todo.ClearDone(db, todosList, todo.RenderTodos, callClear)
			return flags, nil
		case "-clear":
			flags.itemFlags.Parse(commands[1:])
		case "ls":
			todo.RenderTodos(db, callClear)
			return flags, nil
		default:
			flags.itemFlags.Parse(commands[2:])
		}
	}
	return flags, nil
}

func run(db *db.AswanDB) (bool, error) {
	//Initial State
	todosList, err := todo.NewTodoList().Populate(db)
	if err != nil {
		fmt.Println("\nfailed to get todos")
		return false, err
	}
	//Flag Decleration
	flags, err := flagService(os.Args, db, todosList)
	if err != nil {
		fmt.Println("\nFailed to init flags")
		return false, err
	}
	//Handlers
	if *flags.tickFlag {
		possibleInt, err := strconv.ParseInt(flags.nameArg, 10, 64)
		if err != nil {
			possibleInt = -1
		}

		i := slices.IndexFunc(*todosList, func(t *todo.Todo) bool {
			return t.Desc == flags.nameArg
		})
		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return true, nil
		}
		if i == -1 && possibleInt != -1 {
			(*todosList)[possibleInt].ChangeDone(db)
		}
		if possibleInt == -1 && i != -1 {
			(*todosList)[i].ChangeDone(db)
		}
		todosList, err = todo.RenderTodos(db, callClear)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.newFlag {
		i := slices.IndexFunc(*todosList, func(t *todo.Todo) bool {
			return t.Desc == flags.nameArg
		})
		if i != -1 {
			fmt.Printf("\nItem already exists: %s\n", flags.nameArg)
			return true, nil
		}
		ni := todo.NewTodo(flags.nameArg)
		ni.Create(db)
		_, err = todo.RenderTodos(db, callClear)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.deleteFlag {
		possibleInt, err := strconv.ParseInt(flags.nameArg, 10, 64)
		if err != nil {
			possibleInt = -1
		}
		i := slices.IndexFunc(*todosList, func(t *todo.Todo) bool {
			return t.Desc == flags.nameArg
		})

		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return true, nil
		}

		if i == -1 && possibleInt != -1 {
			(*todosList)[possibleInt].Delete(db)
		}

		if possibleInt == -1 && i != -1 {
			(*todosList)[i].Delete(db)
		}
		_, err = todo.RenderTodos(db, callClear)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.clearFlag {
		for _, td := range *todosList {
			td.Delete(db)
		}
		_, err = todo.RenderTodos(db, callClear)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
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
