package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGen(t *testing.T) {
	assert := assert.New(t)

	testCases := []struct {
		src         string
		expected    string
		typeMapping map[string]*Type
	}{
		{
			src: `package main

type Type struct {
	ID int64
}

// TypeDAO implements DAO for Type
type TypeDAO struct {}

// zeroType is the zero value of Type
var zeroType Type

// NewTypeDAO returns a new TypeDAO
func NewTypeDAO() *TypeDAO {
	return &TypeDAO{}
}

// Get loads the Type with the given id.
func (dao *TypeDAO) Get(id int64) (*Type, error) {
	var m Type
	m.ID = id
	return &m, nil
}

// Empty is empty.
func (dao *TypeDAO) Empty() {
}
`,
			expected: `package main

// ConcreteDAO implements DAO for Concrete
type ConcreteDAO struct {
}

// zeroConcrete is the zero value of Concrete
var zeroConcrete Concrete

// NewConcreteDAO returns a new ConcreteDAO
func NewConcreteDAO() *ConcreteDAO {
	return &ConcreteDAO{}
}

// Get loads the Concrete with the given id.
func (dao *ConcreteDAO) Get(id int64) (*Concrete, error) {
	var m Concrete
	m.ID = id
	return &m, nil
}

// Empty is empty.
func (dao *ConcreteDAO) Empty() {
}
`,
			typeMapping: map[string]*Type{
				"Type": {
					Name: "Concrete",
				},
			},
		},
		{
			src: `package main

type Type struct {
	ID int64
}

type TypeDAO struct {}

var zero Type

func NewTypeDAO() *TypeDAO {
	return &TypeDAO{}
}

func (dao *TypeDAO) Get(id int64) (*Type, error) {
	var m Type
	m.ID = id
	return &m, nil
}

func (dao *TypeDAO) Empty() {
}
`,
			expected: `package main

type ConcreteDAO struct {
}

var zero Concrete

func NewConcreteDAO() *ConcreteDAO {
	return &ConcreteDAO{}
}
func (dao *ConcreteDAO) Get(id int64) (*Concrete, error) {
	var m Concrete
	m.ID = id
	return &m, nil
}
func (dao *ConcreteDAO) Empty() {
}
`,
			typeMapping: map[string]*Type{
				"Type": {
					Name: "Concrete",
				},
			},
		},
		{
			src: `package main

type Number int

const ZeroNumber Number = 0

var SomeNumber Number = 42

func AddNumber(a, b Number) Number {
	return a + b
}

func SubNumber(a, b Number) Number {
	return a - b
}
`,
			expected: `package main

const ZeroInt64 int64 = 0

var SomeInt64 int64 = 42

func AddInt64(a, b int64) int64 {
	return a + b
}
func SubInt64(a, b int64) int64 {
	return a - b
}
`,
			typeMapping: map[string]*Type{
				"Number": {
					Name: "int64",
				},
			},
		},
		{
			src: `package foo

type Type struct {
	ID int64
}

func FooType(a Type) Type {
	baz(a.ID)
	return barType(a)
}

func barType(a Type) Type {
	a.ID = 42
}

func baz(id int64) {
	id = 0
}
`,
			expected: `package foo

func FooConcrete(a Concrete) Concrete {
	baz(a.ID)
	return barConcrete(a)
}
func barConcrete(a Concrete) Concrete {
	a.ID = 42
}
`,
			typeMapping: map[string]*Type{
				"Type": {
					Name: "Concrete",
				},
			},
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			inBuff := bytes.NewBufferString(tc.src)
			outBuff := &bytes.Buffer{}
			gen(inBuff, "in.go", "", tc.typeMapping, outBuff, "out.go")

			assert.Equal(tc.expected, outBuff.String())
		})
	}
}
