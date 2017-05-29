package main

import "fmt"

//go:generate rei -in=type.go -out=concrete.go "Type=github.com/nkovacs/rei/example/models.Concrete"
//go:generate rei -in=type.go -out=concrete2gen.go "Type=Concrete2"

func main() {
	dao := NewTypeDAO()
	fmt.Println(dao.Get(12))
	daoC := NewConcreteDAO()
	fmt.Println(daoC.Get(12))
	daoC2 := NewConcrete2DAO()
	fmt.Println(daoC2.Get(12))
}
