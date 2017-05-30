package main

import "fmt"

//go:generate rei -in=slicemap.go -out=stringintmap.go "KeyType=string,ValueType=int"

func main() {
	var m StringIntSliceMap = map[string][]int{
		"one": {1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		"two": {11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
	}
	fmt.Println(m.Flatten())
}
