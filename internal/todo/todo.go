package todo

import (
	"fmt"
	"os"
	"time"

	"github.com/brettearle/aswan/internal/db"
)

const (
	DONE     = "✅"
	NOT_DONE = "❌"
)

type Todolist []*Todo

func NewTodoList() *Todolist {
	return &Todolist{}
}

func (ls *Todolist) Populate(db *db.AswanDB) (*Todolist, error) {
	rows, err := db.GetAllTodos()
	if err != nil {
		fmt.Printf("error: %v", err)
	}
	var res Todolist
	for rows.Next() {
		var item Todo
		if err := rows.Scan(&item.ID, &item.Desc, &item.Done, &item.DoneTime, &item.Board); err != nil {
			fmt.Printf("Scan Error: %v\n", err)
			return nil, err
		}
		res = append(res, &item)
	}
	if len(res) >= 0 {
		return &res, nil
	} else {
		panic("res has no length in populate")
	}
}

type Todo struct {
	ID       int
	Done     bool
	Desc     string
	DoneTime string
	Board    string
}

func (t Todo) String() string {
	return fmt.Sprintf("{desc: %v, done: %v } \n", t.Desc, t.Done)
}

func NewTodo(desc string) *Todo {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("No WD: %v", err)
	}

	i := &Todo{
		Done:     false,
		Desc:     desc,
		DoneTime: time.Now().Format(time.RFC822),
		Board:    wd,
	}
	return i
}

func (i *Todo) Create(db *db.AswanDB) (bool, error) {
	res, err := db.CreateTodo(i.Desc, i.Done, i.DoneTime, i.Board)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	i.ID = int(id)
	return true, err
}

func (i *Todo) Delete(db *db.AswanDB) (bool, error) {
	_, err := db.DeleteTodo(i.ID)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	fmt.Printf("\nDeleted: %+v\n", i)
	return true, nil
}

func (i *Todo) ChangeDone(db *db.AswanDB) (bool, error) {
	if i.Done {
		i.Done = false
	} else {
		i.Done = true
		i.DoneTime = time.Now().Format(time.RFC822)
	}
	_, err := db.UpdateTodo(i.ID, i.Desc, i.Done, i.DoneTime, i.Board)
	if err != nil {
		fmt.Printf("error: %v", err)
		return false, err
	}
	return true, nil
}

type clearTerminal func()

// List all and Print
type RenderList func(db *db.AswanDB, clearTerm clearTerminal, listAll bool) (ls *Todolist, err error)

func RenderTodos(db *db.AswanDB, clearTerm clearTerminal, listAll bool) (*Todolist, error) {
	clearTerm()
	ls, err := NewTodoList().Populate(db)
	if err != nil {
		fmt.Println("failed to get list")
		return nil, err
	}
	if len(*ls) == 0 {
		fmt.Println("No items retrieved")
		return ls, nil
	}
	fmt.Println("")
	if listAll {
		fmt.Println("")
	}
	wd, _ := os.Getwd()
	var count int8 = 0
	for i, todo := range *ls {
		if !listAll {
			if todo.Board != wd {
				continue
			}
			count++
			if count == 1 {
				fmt.Printf("\n===== BOARD: %v =====\n", todo.Board)
			}
		}
		if todo.Done {
			fmt.Printf("%s %d: %s %v \n", DONE, i, todo.Desc, todo.DoneTime)
		} else {
			fmt.Printf("%s %d: %s \n", NOT_DONE, i, todo.Desc)
		}
	}
	if count == 0 && !listAll {
		fmt.Println("Nothing for this directory")
		//   TODO:  fmt.Println("Use the render all command")
		return ls, nil
	}
	return ls, nil
}

// Clear Done
func ClearDone(
	db *db.AswanDB,
	ls *Todolist,
	render RenderList,
	clearTerm clearTerminal,
) (bool, error) {
	for _, item := range *ls {
		if item.Done {
			_, err := item.Delete(db)
			if err != nil {
				fmt.Printf("failed to delete item: %+v.\n Error: %v", item, err)
				return false, err
			}
		}
	}
	render(db, clearTerm, false)
	return true, nil
}
