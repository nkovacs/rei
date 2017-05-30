package main

type Type struct {
	ID int64
}

func FooType(a Type) Type {
	baz(a.ID)
	return barType(a)
}

func barType(a Type) Type {
	a.ID = 42
	return a
}

func baz(id int64) {
	id = 0
}
