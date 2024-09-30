package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/brettearle/aswan/internal/db"
)

const (
	DONE     = "✅"
	NOT_DONE = "❌"
)

type success bool

type todoList []*todo

func newTodoList() *todoList {
	return &todoList{}
}

func (ls *todoList) populate(db *db.AswanDB) (*todoList, error) {
	rows, err := db.GetAllTodos()
	if err != nil {

	}
	var res todoList
	for rows.Next() {
		var item todo
		if err := rows.Scan(&item.id, &item.desc, &item.done); err != nil {
			fmt.Printf("Scan Error: %v\n", err)
			return nil, err
		}
		res = append(res, &item)
	}
	return &res, nil
}

type todo struct {
	id   int
	done bool
	desc string
}

func (t todo) String() string {
	return fmt.Sprintf("{desc: %v, done: %v } \n", t.desc, t.done)
}

func newTodo(desc string) *todo {
	i := &todo{
		done: false,
		desc: desc,
	}
	return i
}

func (i *todo) create(db *db.AswanDB) (success, error) {
	res, err := db.CreateTodo(i.desc, i.done)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	i.id = int(id)
	return true, err
}

func (i *todo) delete(db *db.AswanDB) (success, error) {
	_, err := db.DeleteTodo(i.id)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	fmt.Printf("\nDeleted: %+v\n", i)
	return true, nil
}

func (i *todo) tickUntick(db *db.AswanDB) (success, error) {
	if i.done {
		i.done = false
	} else {
		i.done = true
	}
	_, err := db.UpdateTodo(i.id, i.desc, i.done)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	return true, nil
}

// List all and Print
type RenderList func(db *db.AswanDB) (ls *todoList, err error)

func renderList(db *db.AswanDB) (*todoList, error) {
	callClear()
	ls, err := newTodoList().populate(db)
	if err != nil {
		fmt.Println("failed to get list")
		return nil, err
	}
	if len(*ls) == 0 {
		fmt.Println("No todos")
		return ls, nil
	}
	fmt.Println("")
	for i, todo := range *ls {
		if todo.done {
			fmt.Printf("%s %d: %s \n", DONE, i, todo.desc)
		} else {
			fmt.Printf("%s %d: %s \n", NOT_DONE, i, todo.desc)
		}
	}
	return ls, nil
}

// Clear Done
func clearDone(db *db.AswanDB, ls *todoList, render RenderList) (success, error) {
	for _, item := range *ls {
		if item.done {
			_, err := item.delete(db)
			if err != nil {
				fmt.Printf("failed to delete item: %+v.\n Error: %v", item, err)
				return false, err
			}
		}
	}
	render(db)
	return true, nil
}

type todoFlags struct {
	nameArg string
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
		nameArg: desc,
	}
}

func flagService(
	args []string,
	db *db.AswanDB,
	todosList *todoList,
) (*todoFlags, error) {

	var err error
	//Commands
	commands := args
	if len(commands) == 1 {
		renderList(db)
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
			clearDone(db, todosList, renderList)
			return flags, nil
		case "-clear":
			flags.itemFlags.Parse(commands[1:])
		case "ls":
			renderList(db)
			return flags, nil
		default:
			flags.itemFlags.Parse(commands[2:])
		}
	}
	return flags, nil
}

func run(db *db.AswanDB) (success, error) {
	//Initial State
	todosList, err := newTodoList().populate(db)
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

		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == flags.nameArg
		})
		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return true, nil
		}
		if i == -1 && possibleInt != -1 {
			(*todosList)[possibleInt].tickUntick(db)
		}
		if possibleInt == -1 && i != -1 {
			(*todosList)[i].tickUntick(db)
		}
		todosList, err = renderList(db)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.newFlag {
		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == flags.nameArg
		})
		if i != -1 {
			fmt.Printf("\nItem already exists: %s\n", flags.nameArg)
			return true, nil
		}
		ni := newTodo(flags.nameArg)
		ni.create(db)
		_, err = renderList(db)
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
		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == flags.nameArg
		})

		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return true, nil
		}

		if i == -1 && possibleInt != -1 {
			(*todosList)[possibleInt].delete(db)
		}

		if possibleInt == -1 && i != -1 {
			(*todosList)[i].delete(db)
		}
		_, err = renderList(db)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.clearFlag {
		for _, td := range *todosList {
			td.delete(db)
		}
		_, err = renderList(db)
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
