package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

const (
	DONE     = "✅"
	NOT_DONE = "❌"
)

type success bool

type todoList []*todo

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

func (i *todo) create(db *aswanDB) (success, error) {
	res, err := db.createTodo(i)
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

func (i *todo) delete(db *aswanDB) (success, error) {
	_, err := db.deleteTodo(i.id)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	fmt.Printf("\nDeleted: %+v\n", i)
	return true, nil
}

func (i *todo) tickUntick(db *aswanDB) (success, error) {
	if i.done {
		i.done = false
	} else {
		i.done = true
	}
	_, err := db.updateTodo(i)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	return true, nil
}

// List all and Print
type RenderList func(db *aswanDB) (ls *todoList, err error)

func renderList(db *aswanDB) (*todoList, error) {
	ls, err := db.getAllTodos()
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
func clearDone(db *aswanDB, ls *todoList, render RenderList) (success, error) {
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
	DB *aswanDB,
	todosList *todoList,
) (*todoFlags, error) {

	var err error
	//Commands
	commands := args
	if len(commands) == 1 {
		renderList(DB)
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
			fmt.Printf("\n%s\n", DB.path)
			return flags, nil
		case "timer":
			fmt.Println("\ntimer not yet implemented")
			return flags, nil
		case "rmDone":
			clearDone(DB, todosList, renderList)
			return flags, nil
		case "-clear":
			flags.itemFlags.Parse(commands[1:])
		case "ls":
			renderList(DB)
			return flags, nil
		default:
			flags.itemFlags.Parse(commands[2:])
		}
	}
	return flags, nil
}

func run(DB *aswanDB) (success, error) {
	//Initial State
	todosList, err := DB.getAllTodos()
	if err != nil {
		fmt.Println("\nfailed to get todos")
		return false, err
	}
	//Flag Decleration
	flags, err := flagService(os.Args, DB, todosList)
	if err != nil {
		fmt.Println("\nFailed to init flags")
		return false, err
	}
	//-- End Flag Decleration --

	//Exploration
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
			(*todosList)[possibleInt].tickUntick(DB)
		}
		if possibleInt == -1 && i != -1 {
			(*todosList)[i].tickUntick(DB)
		}
		todosList, err = renderList(DB)
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
		ni.create(DB)
		_, err = renderList(DB)
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
			(*todosList)[possibleInt].delete(DB)
		}

		if possibleInt == -1 && i != -1 {
			(*todosList)[i].delete(DB)
		}
		_, err = renderList(DB)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	if *flags.clearFlag {
		for _, td := range *todosList {
			td.delete(DB)
		}
		_, err = renderList(DB)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return false, err
		}
	}

	return true, nil
}

func main() {
	//DB Initialization
	DB, err := dbInit(getDBPath())
	if err != nil {
		panic("no DB able to be initialized")
	}
	defer DB.db.Close()

	//RUN RUN RUN 
	_, err = run(DB)
	if err != nil {
		fmt.Printf("\nRun failed with: %v", err)
	}
}
