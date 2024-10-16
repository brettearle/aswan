package flagger

import (
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/brettearle/aswan/internal/db"
	"github.com/brettearle/aswan/internal/terminal"
	"github.com/brettearle/aswan/internal/timer"
	"github.com/brettearle/aswan/internal/todo"
)

type TodoFlags struct {
	Name      string
	ItemFlags *flag.FlagSet
	New       *bool
	Tick      *bool
	Delete    *bool
	Clear     *bool
}

func newTodoFlags(desc string) *TodoFlags {
	iF := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	return &TodoFlags{
		ItemFlags: iF,
		New:       iF.Bool("n", false, "New todo"),
		Tick:      iF.Bool("t", false, "Completes a todo"),
		Delete:    iF.Bool("d", false, "Deletes a todo"),
		Clear:     iF.Bool("clear", false, "Deletes all todos"),
		Name:      desc,
	}
}

func FlagService(
	args []string,
	db *db.AswanDB,
	todosList *todo.Todolist,
) (*TodoFlags, error) {

	var err error
	//Commands
	commands := args
	if len(commands) == 1 {
		todo.RenderTodos(db, terminal.CallClear, false)
		return newTodoFlags(""), err
	}

	flags := newTodoFlags(args[1])

	_, clear, flagFirst := strings.Cut(commands[1], "-")

	//Handles structure
	if flagFirst && clear != "clear" {
		fmt.Println("\nPlease provide item string first")
		fmt.Println("example: `CMD> aswan 'im a item' -n -t`")
		flags.ItemFlags.PrintDefaults()
		return flags, err
	}

	if len(commands) > 1 {
		switch commands[1] {
		//cases for commands go here
		case "help":
			fmt.Println("")
			fmt.Println("ls - list all todos for current board, this is default command")
			fmt.Println("")
			fmt.Println("lsa - list all todos in the database")
			fmt.Println("")
			fmt.Println("dbPath - Get path to local db")
			fmt.Println("")
			fmt.Println("rmDone - Remove completed todos")
			fmt.Println("")
			flags.ItemFlags.PrintDefaults()
			return flags, nil
		case "dbPath":
			fmt.Printf("\n%s\n", db.Path)
			return flags, nil
		case "timer":
			t := timer.NewTermPomoTimer(10, 10, 10, 10)
			//Concurrent here maybe???
			t.Start()
			return flags, nil
		case "rmDone":
			todo.ClearDone(db, todosList, todo.RenderTodos, terminal.CallClear)
			return flags, nil
		case "-clear":
			flags.ItemFlags.Parse(commands[1:])
		case "ls":
			todo.RenderTodos(db, terminal.CallClear, false)
			return flags, nil
		case "lsa":
			todo.RenderTodos(db, terminal.CallClear, true)
			return flags, nil
		default:
			if len(commands) > 1 {
				flags.ItemFlags.Parse(commands[2:])
			}
			if len(commands) == 2 {
				commands = append(commands, "-t")
				flags.ItemFlags.Parse(commands[2:])
			}
		}
	}
	return flags, nil
}

func TickHandler(list *todo.Todolist, flags *TodoFlags, db *db.AswanDB) (*todo.Todolist, error) {
	possibleInt, err := strconv.ParseInt(flags.Name, 10, 64)
	if err != nil {
		possibleInt = -1
	}

	i := slices.IndexFunc(*list, func(t *todo.Todo) bool {
		return t.Desc == flags.Name
	})
	if i == -1 && possibleInt == -1 {
		fmt.Println("\nNo item by that name")
		return list, nil
	}
	if i == -1 && possibleInt != -1 {
		(*list)[possibleInt].ChangeDone(db)
	}
	if possibleInt == -1 && i != -1 {
		(*list)[i].ChangeDone(db)
	}
	list, err = todo.RenderTodos(db, terminal.CallClear, false)
	if err != nil {
		fmt.Println("\nCouldn't get updated list")
		return list, err
	}
	return list, nil
}

func NewHandler(list *todo.Todolist, flags *TodoFlags, db *db.AswanDB) (*todo.Todolist, error) {
	i := slices.IndexFunc(*list, func(t *todo.Todo) bool {
		return t.Desc == flags.Name
	})
	if i != -1 {
		fmt.Printf("\nItem already exists: %s\n", flags.Name)
		return list, nil
	}
	ni := todo.NewTodo(flags.Name)
	ni.Create(db)
	var err error
	list, err = todo.RenderTodos(db, terminal.CallClear, false)
	if err != nil {
		fmt.Println("\nCouldn't get updated list")
		return list, err
	}
	return list, nil
}

func DeleteHandler(list *todo.Todolist, flags *TodoFlags, db *db.AswanDB) (*todo.Todolist, error) {
	possibleInt, err := strconv.ParseInt(flags.Name, 10, 64)
	if err != nil {
		possibleInt = -1
	}
	i := slices.IndexFunc(*list, func(t *todo.Todo) bool {
		return t.Desc == flags.Name
	})

	if i == -1 && possibleInt == -1 {
		fmt.Println("\nNo item by that name")
		return list, nil
	}

	if i == -1 && possibleInt != -1 {
		(*list)[possibleInt].Delete(db)
	}

	if possibleInt == -1 && i != -1 {
		(*list)[i].Delete(db)
	}
	list, err = todo.RenderTodos(db, terminal.CallClear, false)
	if err != nil {
		fmt.Println("\nCouldn't get updated list")
		return list, err
	}
	return list, nil
}

func ClearHandler(list *todo.Todolist, flags *TodoFlags, db *db.AswanDB) (*todo.Todolist, error) {
	for _, td := range *list {
		td.Delete(db)
	}
	var err error
	list, err = todo.RenderTodos(db, terminal.CallClear, false)
	if err != nil {
		fmt.Println("\nCouldn't get updated list")
		return list, err
	}
	return list, nil
}
