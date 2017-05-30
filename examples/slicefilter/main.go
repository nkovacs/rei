package main

import "fmt"

//go:generate rei -in=slicefilter.go -out=intfilter.go "Type=int"

func main() {
	var s IntSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	evens := s.Where(func(val int) bool { return val%2 == 0 })
	fmt.Printf("Evens: %+v\n", evens)
}
