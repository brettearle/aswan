package main

import "testing"

func TestTodo(t *testing.T) {
	testDB, err := dbInit(`:memory:`)
	if err != nil {
		t.Errorf("failed to init in mem test db: %s", err)
	}

	t.Run("new todo has correct desc", func(t *testing.T) {
		got := newTodo("test")
		want := "test"
		if got.desc != want {
			t.Errorf("got %s want %s", got.desc, want)
		}
	})
	t.Run("init with done equal to false", func(t *testing.T) {
		got := newTodo("test")
		want := false
		if got.done != want {
			t.Errorf("got %v want %v", got.done, want)
		}
	})
	t.Run("create a todo", func(t *testing.T) {
		got := newTodo("test")
		got.create(testDB)
		want := 1
		if got.id != want {
			t.Errorf("got %v want %v", got.id, want)
		}
	})
	t.Run("tickUntick should swap value of Done on todo", func(t *testing.T) {
		got := newTodo("test")
		got.tickUntick(testDB)
		want := true
		if got.done != want {
			t.Errorf("got %v want %v", got.done, want)
		}
	})
	t.Run("delete should return true when successful", func(t *testing.T) {
		td := newTodo("test")
		got, _ := td.delete(testDB)
		var want success = true
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("render list should return the correct list", func(t *testing.T) {
		_, err := newTodo("test").create(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		ls, _ := renderList(testDB)
		got := *ls
		want := todo{
			id:   1,
			desc: "test",
			done: false,
		}
		if *got[0] != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
	t.Run("should only have todos with false left in list after clearDone", func(t *testing.T) {
		testTD := newTodo("test")
		testTD1 := newTodo("test1")
		testDB.createTodo(testTD.desc, testTD.done)
		testDB.createTodo(testTD1.desc, testTD1.done)
		ls, err := newTodoList().populate(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		list := *ls
		list[0].tickUntick(testDB)
		clearDone(testDB, &list, renderList)
		clLs, err := newTodoList().populate(testDB)
		if err != nil {
			t.Errorf("failed new todo %v", err)
		}
		clearedList := *clLs
		got := clearedList[0].done
		want := false
		if got != want {
			t.Errorf("got %v want %v", got, want)
		}
	})
}
