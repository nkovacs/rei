package main

import "fmt"

//go:generate rei -in=numbers.go -out=int64.go "Number=int64"

func main() {
	fmt.Println(AddInt64(1, 1))
	fmt.Println(SubInt64(2, 1))
}
