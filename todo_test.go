package main

import "testing"

func TestTodo(t *testing.T) {
	testDB, err := dbInit(`:memory:`)
	if err != nil {
		t.Errorf("failed to init in mem test db: %s", err)
	}
	defer testDB.db.Close()

	t.Run("new todo has correct desc", func(t *testing.T) {
		got := newTodo("test")
		want := "test"
		if got.desc != want {
			t.Errorf("got %s want %s", got.desc, want)
		}
	})
	t.Run("init with false as done", func(t *testing.T) {
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

}