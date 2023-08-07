package main

import "fmt"

type h struct {
	a int
}

func main() {
	a := (*h)(nil)
	a = &h{}
	fmt.Println(a.a)
}
