package main

import (
	"flag"
	"fmt"
)


type item struct{
	done bool
	desc string
}
func newItem(desc string)(*item, error){
	i := &item{
		done: false,
		desc: desc,
	}
	return i, nil
}

func main() {
	//Flag Decleration
	newFlag := flag.String("n", "Item Description", "New Item: Value is Item Description")
	flag.Parse()

	//Exploration
	fmt.Printf("newFlag: %s", *newFlag )
	i,err := newItem("tested")
	if (err != nil) {
		fmt.Println("Item broke")
	}
	fmt.Printf("item: %v", i)
}

