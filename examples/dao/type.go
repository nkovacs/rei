package main

import "fmt"

type Type struct {
	ID int64
}

// TypeDAO is a data access object for Type
type TypeDAO struct{}

// NewTypeDAO creates a new TypeDAO
func NewTypeDAO() *TypeDAO {
	return &TypeDAO{}
}

// Get loads the Type with the given id.
func (dao *TypeDAO) Get(id int64) (*Type, error) {
	var m Type
	m.ID = id
	return &m, nil
}

// Set sets the Type with the given id.
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
	if c := 0; true {
		fmt.Println("something", c)
	} else if false {
		fmt.Println("something else")
	} else {
		fmt.Println("no way")
	}
}
