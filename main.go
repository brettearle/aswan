package main

import "fmt"


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
	i,err := newItem("tested")
	if (err != nil) {
		fmt.Println("Item broke")
	}
	fmt.Printf("item: %v", i)
}

