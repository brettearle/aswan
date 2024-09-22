package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type itemList []*item

type item struct {
	id int
	done bool
	desc string
}

func (t item) String() string {
	return fmt.Sprintf("\n{%v %v}\n", t.done, t.desc)
}

// helpers
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
	fmt.Printf("Created: %+v\n with res: %v\n", i, res)
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
	res, err := db.updateTodo(i)	
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	fmt.Printf("Updated: %+v\n with res: %v\n", i, res)
}

func main() {
	//DB Initialization
	testDB, err := dbInit("./test.db")
	if err != nil {
		fmt.Println("failed to init DB")
	}
	//Flag Decleration
	itemFlags := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	newFlag := itemFlags.Bool("n", false, "New Item")
	tickFlag := itemFlags.Bool("t", false, "Completes an item")
	listFlag := itemFlags.Bool("ls", false, "List items")
	//Commands
	commands := os.Args
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
		i := newItem(itemDesc)
		i.create(testDB)
		i.tickUntick(testDB)
	}
	if *newFlag {
		ni := newItem(itemDesc)
		ni.create(testDB)
		ni.tickUntick(testDB)
	}
	if *listFlag {
		res, err := testDB.getAllTodos()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Println("All Items -")
		for _, val := range *res {
			fmt.Printf("%v: %v\n", val.id, val)
			// val.delete(testDB)
		}
	}
	fmt.Printf("DB: %+v\n", testDB)
}

