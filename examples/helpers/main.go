package main

import "fmt"

//go:generate rei -in=type.go -out=concrete.go "Type=Concrete"

type Concrete struct {
	ID int64
}

func main() {
	var c Concrete
	c = FooConcrete(c)
	fmt.Println(c.ID)
}
