package main

import "fmt"

type Type struct {
	ID int64
}

type TypeDAO struct{}

type TypeDependant struct {
	m Type
}

func NewTypeDAO() *TypeDAO {
	return &TypeDAO{}
}

func (dao *TypeDAO) Get(id int64) (*Type, error) {
	var m Type
	m.ID = id
	return &m, nil
}

func (dao *TypeDAO) Set(id int64) {
	type (
		Type struct {
			IDD int64
		}
	)

	var m Type
	m.IDD = id
}

func (dao *TypeDAO) GetType() {
}

func SomethingType(d TypeDependant) {
}

func commonhelper() {
}

func (dao *TypeDAO) Foo() {
	commonhelper()
}

var zeroType Type

var TypeA, TypeB Type

func foobarTypeA() Type {
	return TypeA
}

func foobarTypeB() Type {
	return TypeB
}

func foobarTypeBoth() (Type, Type) {
	return TypeA, TypeB
}

func foobarType(t Type) {
}

func barbazType(d TypeDependant) {
}

func barfooType() {
	foobarType(zeroType)
}

func someHelper() {
	fmt.Println("just a helper, not copied")
}
