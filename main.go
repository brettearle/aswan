package main

import (
	"flag"
	"fmt"
	"os"
)

type ItemList []Item

type Item struct {
	done bool
	desc string
}

func NewItem(desc string) *Item {
	i := &Item{
		done: false,
		desc: desc,
	}
	return i
}

func (i *Item) Create() {
	fmt.Printf("Created: %+v\n", i)
}

func (i *Item) TickUntick() {
	if i.done {
		i.done = false
	}
	i.done = true
	fmt.Printf("Ticked: %+v\n", i)
}

func (i *Item) Print() {
	fmt.Printf("Current Item: %+v\n", i)
}

func main() {
	//Flag Decleration
	itemFlags := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	newFlag := itemFlags.Bool("n", false, "New Item")
	tickFlag := itemFlags.Bool("t", false, "completes an item")
	listFlag := itemFlags.Bool("ls", false, "list all items")
	fmt.Println(*listFlag)
	//Commands
	commands := os.Args
	//TODO sort out if just flags passed in
	if len(commands) > 1 {
		switch commands[1] {
		//cases for commands go here 
		default:
			itemFlags.Parse(commands[2:])
		}
	} else {
		fmt.Println("Need this structure for aswan command")
		fmt.Println("aswan 'item name' -flag -fl.....")
		return
	}
	//-- End Flag Decleration --
	//Arguments
	itemArg := commands[1]
	//-- End Args --

	//Exploration
	i := NewItem(itemArg)
	if *tickFlag {
		i.TickUntick()
		i.Print()
	}
	if *newFlag {
		ni := NewItem(itemArg)
		ni.Create()
		ni.Print()
	}
}
