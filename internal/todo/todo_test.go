package todo

import (
	"reflect"
	"testing"

	"github.com/brettearle/aswan/internal/db"
)

func TestTodo(t *testing.T) {
	callClear := func() {}
	testDB, err := db.DbInit(`:memory:`)
	if err != nil {
		t.Errorf("failed to init in mem test db: %s", err)
	}

	t.Run("new todo has correct desc", func(t *testing.T) {
		got := NewTodo("test")
		want := "test"
		if got.Desc != want {
			t.Errorf("got %s want %s", got.Desc, want)
		}
	})
	t.Run("init with done equal to false", func(t *testing.T) {
		got := NewTodo("test")
		want := false
		if got.Done != want {
			t.Errorf("got %v want %v", got.Done, want)
		}
	})
	t.Run("create a todo", func(t *testing.T) {
		got := NewTodo("test")
		got.Create(testDB)
		want := 1
		if got.ID != want {
			t.Errorf("got %v want %v", got.ID, want)
		}
		got.Delete(testDB)
	})
	t.Run("tickUntick should swap value of Done on todo", func(t *testing.T) {
		got := NewTodo("test")
		got.ChangeDone(testDB)
		want := true
		if got.Done != want {
			t.Errorf("got %v want %v", got.Done, want)
		}
	})
	t.Run("delete should return true when successful", func(t *testing.T) {
		td := NewTodo("test")
		got, _ := td.Delete(testDB)
		var want = true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("render list should return the correct list", func(t *testing.T) {
		_, err := NewTodo("test").Create(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		ls, _ := RenderTodos(testDB, callClear)
		got := ls
		want, _ := NewTodoList().Populate(testDB)
		if !reflect.DeepEqual(*got, *want) {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should only have todos with false left in list after clearDone", func(t *testing.T) {
		testTD := NewTodo("test")
		testTD1 := NewTodo("test1")
		testDB.CreateTodo(testTD.Desc, testTD.Done, testTD.DoneTime, testTD.Board)
		testDB.CreateTodo(testTD1.Desc, testTD1.Done, testTD1.DoneTime, testTD1.Board)
		ls, err := NewTodoList().Populate(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		list := *ls
		list[0].ChangeDone(testDB)
		ClearDone(testDB, &list, RenderTodos, callClear)
		clLs, err := NewTodoList().Populate(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		clearedList := *clLs
		got := clearedList[0].Done
		want := false
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
