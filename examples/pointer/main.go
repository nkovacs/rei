package main

import (
	"fmt"
	"os"

	test "github.com/nkovacs/rei/examples/pointer/go-test"
)

//go:generate rei -in=reader.go -out=readerfile.go "Reader=*os.File"
//go:generate rei -in=reader.go -out=readertest.go "Reader=*(\"github.com/nkovacs/rei/examples/pointer/go-test\")test.TestReader"

func main() {
	f, err := os.Open("main.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s, err := ReadAllStringFromReader(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", s)

	s, err = ReadAllStringFromTestReader(&test.TestReader{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", s)
}
