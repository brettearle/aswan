package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type ItemList []*Item

type Item struct {
	done bool
	desc string
}

func (t Item) String() string {
	return fmt.Sprintf("\n{%v %v}\n", t.done, t.desc)
}

// helpers
func NewItem(desc string) *Item {
	i := &Item{
		done: false,
		desc: desc,
	}
	return i
}
func GetItem(desc string, iL ItemList) (*Item, error) {
	for _, item := range iL {
		if item.desc == desc {
			return item, nil
		}
	}
	return nil, errors.New("Item Not In List")
}

func (i *Item) Create() {
	fmt.Printf("Created: %+v\n", i)
}

func (i *Item) TickUntick() {
	if i.done {
		i.done = false
	} else {
		i.done = true
	}
	fmt.Printf("Ticked: %+v\n", i)
}

func (i *Item) Print() {
	fmt.Printf("Current Item: %+v\n", i)
}

func main() {
	//Flag Decleration
	itemFlags := flag.NewFlagSet("Todo items", flag.ContinueOnError)
	newFlag := itemFlags.Bool("n", false, "New Item")
	tickFlag := itemFlags.Bool("t", false, "Completes an item")
	listFlag := itemFlags.Bool("ls", false, "List items")
	fmt.Println(*listFlag)
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
	var testList ItemList
	testList = append(testList, NewItem(itemDesc))
	if *tickFlag {
		i, err := GetItem(itemDesc, testList)
		if err != nil {
			i = NewItem(itemDesc)
			testList = append(testList, i)
		}
		i.TickUntick()
		i.Print()
	}
	if *newFlag {
		ni, err := GetItem(itemDesc, testList)
		if err != nil {
			ni = NewItem(itemDesc)
			testList = append(testList, ni)
		}
		ni.TickUntick()
		ni.Create()
		ni.Print()
	}
	if *listFlag {
		fmt.Println("All Items -")
		fmt.Printf("%v", testList)
	}

}
