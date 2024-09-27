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

func (i *todo) create(db *aswanDB) (bool, error) {
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

func (i *todo) delete(db *aswanDB) (bool, error)  {
	_, err := db.deleteTodo(i.id)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	fmt.Printf("\nDeleted: %+v\n", i)
	return true, nil
}

func (i *todo) tickUntick(db *aswanDB) (bool, error) {
	if i.done {
		i.done = false
	} else {
		i.done = true
	}
	_, err := db.updateTodo(i)
	if err != nil {
		fmt.Printf("error: %v", err);
		return false, err
	}
		return true, nil
}

// List all and Print
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

//Clear Don

func main() {
	//DB Initialization
	dbPath := getDBPath()
	DB, err := dbInit(dbPath)
	if err != nil {
		panic("no DB able to be initialized")
	}
	defer DB.db.Close()

	//Initial State
	todosList, err := DB.getAllTodos()
	if err != nil {
		fmt.Println("\nfailed to get todos")
		return
	}
	//Flag Decleration
	itemFlags := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	newFlag := itemFlags.Bool("n", false, "New todo")
	tickFlag := itemFlags.Bool("t", false, "Completes a todo")
	deleteFlag := itemFlags.Bool("d", false, "Deletes a todo")
	clearFlag := itemFlags.Bool("clear", false, "Deletes all todos")
	//Commands
	commands := os.Args
	if len(commands) == 1 {
		renderList(DB)
		return
	}
	_, clear, flagFirst := strings.Cut(commands[1], "-")

	//Handles structure
	if flagFirst && clear != "clear" {
		fmt.Println("\nPlease provide item string first")
		fmt.Println("example: `CMD> aswan 'im a item' -n -t`")
		itemFlags.PrintDefaults()
		return
	}

	if len(commands) > 1 {
		switch commands[1] {
		//cases for commands go here
		case "help":
			fmt.Println("\nhelp not implemented")
			return
		case "dbPath":
			fmt.Printf("\n%s\n", DB.path)
			return
		case "timer":
			fmt.Println("\ntimer not yet implemented")
			return
		case "-clear":
			itemFlags.Parse(commands[1:])
		case "ls":
			renderList(DB)
			return
		default:
			itemFlags.Parse(commands[2:])
		}
	}
	//-- End Flag Decleration --

	//Arguments
	itemDesc := commands[1]
	//-- End Args --

	//Exploration
	if *tickFlag {
		possibleInt, err := strconv.ParseInt(itemDesc, 10, 64)
		if err != nil {
			possibleInt = -1
		}

		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == itemDesc
		})
		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return
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
			return
		}
	}

	if *newFlag || !*tickFlag && !*deleteFlag && !*newFlag && !*clearFlag {
		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == itemDesc
		})
		if i != -1 {
			fmt.Printf("\nItem already exists: %s\n", itemDesc)
			return
		}
		ni := newTodo(itemDesc)
		ni.create(DB)
		_, err = renderList(DB)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return
		}
	}

	if *deleteFlag {
		possibleInt, err := strconv.ParseInt(itemDesc, 10, 64)
		if err != nil {
			possibleInt = -1
		}
		i := slices.IndexFunc(*todosList, func(t *todo) bool {
			return t.desc == itemDesc
		})

		if i == -1 && possibleInt == -1 {
			fmt.Println("\nNo item by that name")
			return
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
			return
		}
	}

	if *clearFlag {
		for _, td := range *todosList {
			td.delete(DB)
		}
		_, err = renderList(DB)
		if err != nil {
			fmt.Println("\nCouldn't get updated list")
			return
		}
	}
}
