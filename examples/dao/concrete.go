// Code generated by rei. DO NOT EDIT.

package main

import (
	"fmt"

	"github.com/nkovacs/rei/examples/dao/models"
)

// ConcreteDAO is a data access object for Concrete
type ConcreteDAO struct {
}

// NewConcreteDAO creates a new ConcreteDAO
func NewConcreteDAO() *ConcreteDAO {
	return &ConcreteDAO{}
}

// Get loads the Concrete with the given id.
func (dao *ConcreteDAO) Get(id int64) (*models.Concrete, error) {
	var m models.Concrete
	m.ID = id
	return &m, nil
}

// Set sets the Concrete with the given id.
func (dao *ConcreteDAO) Set(id int64) {
	type Type struct {
		IDD int64
	}
	var m Type
	m.IDD = id
}
func (dao *ConcreteDAO) GetType() {
	if c := 0; true {
		fmt.Println("something", c)
	} else if false {
		fmt.Println("something else")
	} else {
		fmt.Println("no way")
	}
}
