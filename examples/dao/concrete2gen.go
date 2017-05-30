package main

import "fmt"

// Concrete2DAO is a data access object for Concrete2
type Concrete2DAO struct {
}

// NewConcrete2DAO creates a new Concrete2DAO
func NewConcrete2DAO() *Concrete2DAO {
	return &Concrete2DAO{}
}

// Get loads the Concrete2 with the given id.
func (dao *Concrete2DAO) Get(id int64) (*Concrete2, error) {
	var m Concrete2
	m.ID = id
	return &m, nil
}

// Set sets the Concrete2 with the given id.
func (dao *Concrete2DAO) Set(id int64) {
	type Type struct {
		IDD int64
	}
	var m Type
	m.IDD = id
}
func (dao *Concrete2DAO) GetType() {
	if c := 0; true {
		fmt.Println("something", c)
	} else if false {
		fmt.Println("something else")
	} else {
		fmt.Println("no way")
	}
}