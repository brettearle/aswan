package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

type itemList []*item

type item struct {
	id   int
	done bool
	desc string
}

func (t item) String() string {
	return fmt.Sprintf("{desc: %v, done: %v } \n", t.desc, t.done)
}

func newItem(desc string) *item {
	i := &item{
		done: false,
		desc: desc,
	}
	return i
}

func (i *item) create(db *aswanDB) {
	res, err := db.createTodo(i)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	i.id = int(id)
}

func (i *item) delete(db *aswanDB) {
	res, err := db.deleteTodo(i.id)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Printf("Deleted: %+v\n with res: %v\n", i, res)
}

func (i *item) tickUntick(db *aswanDB) {
	if i.done {
		i.done = false
	} else {
		i.done = true
	}
	_, err := db.updateTodo(i)
	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

// List all and Print
func renderList(db *aswanDB) (*itemList, error) {
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
		fmt.Printf("%d: %s %t\n", i, todo.desc, todo.done)
	}

	return ls, nil
}

func main() {
	//DB Initialization
	testDB, err := dbInit("./test.db")
	if err != nil {
		fmt.Println("failed to init DB")
	}
	//Initial State and Screen
	todosList, err := renderList(testDB)
	if err != nil {
		fmt.Println("failed to get todos")
		return
	}
	//Flag Decleration
	itemFlags := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	newFlag := itemFlags.Bool("n", false, "New Item")
	tickFlag := itemFlags.Bool("t", false, "Completes an item")
	// listFlag := itemFlags.Bool("ls", false, "List items")
	//Commands
	commands := os.Args
	if len(commands) == 1 {
		return
	}
	_, _, flagFirst := strings.Cut(commands[1], "-")

	//Handles structure
	if flagFirst {
		fmt.Println("Please provide item string first")
		fmt.Println("example: `CMD> aswan 'im a item' -n -t`")
		itemFlags.PrintDefaults()
		return
	}

	if len(commands) > 1 {
		switch commands[1] {
		//cases for commands go here
		case "help":
			fmt.Println("help not implemented")
		case "dbPath":
			fmt.Println("path to DB")
		case "timer":
			fmt.Println("timer not yet implemented")
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
		i := slices.IndexFunc(*todosList, func(t *item) bool {
			return t.desc == itemDesc
		})
		if i == -1 {
			fmt.Println("\nNo item by that name")
		}
		(*todosList)[i].tickUntick(testDB)
		todosList, err = renderList(testDB)
		if err != nil {
			fmt.Println("Couldn't get updated list")
			return
		}
	}
	if *newFlag {
		i := slices.IndexFunc(*todosList, func(t *item) bool {
			return t.desc == itemDesc
		})
		if i != -1 {
			fmt.Printf("\nItem already exists: %s\n", itemDesc)
			return
		}
		ni := newItem(itemDesc)
		ni.create(testDB)
		_, err = renderList(testDB)
		if err != nil {
			fmt.Println("Couldn't get updated list")
			return
		}
	}
	// if *listFlag {
	// 	res, err := testDB.getAllTodos()
	// 	if err != nil {
	// 		fmt.Printf("Error: %v\n", err)
	// 	}
	// 	fmt.Println("All Items -")
	// 	for _, val := range *res {
	// 		fmt.Printf("%v: %v\n", val.id, val)
	// 		// val.delete(testDB)
	// 	}
	// }
	// fmt.Printf("DB: %+v\n", testDB)
}
