package main

import "fmt"

func main() {
	h := test()
	fmt.Println(h)
}

func test() (res int) {
	a := 1
	defer func() { a++; fmt.Println(res) }()
	return a
}
