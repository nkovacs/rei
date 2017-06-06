# rei

Reified generics with code generation for Go

## Installation

```
go get -u github.com/nkovacs/rei
```

## Usage

```
rei -in=generic.go -out=concrete.go 'Type=string'
```

This will replace all instances of `Type` in generic.go with `string`, and write out the result to concrete.go.

The type mapping is specified as `{generic}={concrete},[{generic}={concrete}...]`, where `generic` must be an identifier,
and `concrete` can be one of the following:
- **ConcreteType**: a type in the same package as the generated file, or a builtin type, e.g. `string` (no import will be generated)
- **pkg.ConcreteType**: a type in a different package, will generate `import "pkg"`
- **(pkgPath)pkgAlias.ConcreteType**: a type in a different package, will generate `import pkgAlias "pkgPath"`
- **\*Type**: pointer to another type

### Example

```go
package daos

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
```

Running `rei -in=type.go "Type=models.User"` will generate:

```go
package daos

import "models"

// UserDAO implements DAO for User
type UserDAO struct {
}

// zeroUser is the zero value of User
var zeroUser models.User

// NewUserDAO returns a new UserDAO
func NewUserDAO() *UserDAO {
        return &UserDAO{}
}

// Get loads the User with the given id.
func (dao *UserDAO) Get(id int64) (*models.User, error) {
        var m models.User
        m.ID = id
        return &m, nil
}
```

See the examples directory for more examples.

## How it works

Rei replaces one or more generic types declared in a file with concrete types.
These generic types must be declared in the source file, but they can be anything, including structs.
This allows you to write working, testable generic code that uses methods and struct fields.

It collects all global functions, variables, constants and types that use the generic type,
and generates versions that use the concrete type.
Variables, constants, types, and functions without receivers are renamed by replacing every instance of a
generic type's name with the corresponding concrete type's name, keeping the case of the first character.

E.g.:
- With the mapping `Type=Foo`, `func FrobnizeType` will be renamed to `func FrobnizeFoo`
and `type typeHelper` will be renamed to `type fooHelper`.
- With the mapping `Type=string`, `func FrobnizeType` will be renamed to `func FrobnizeString`
and `type typeHelper` will be renamed to `type stringHelper`.

Methods are not renamed, since the receiver type's name makes them unique.

## Known limitations

- Only accepts a single file as input.
- Dot imports are copied over to the generated file, but they cannot be removed by goimports.
  If the imported package is not used, this will cause a compilation error.
- Declarations that are not generic are not copied to the generated file, since this can cause duplicate declarations.
  E.g. a generic function cannot call a non-generic function when generating into a different directory.
- When generating into a different directory, the generated package name is the directory name.
- If the generic declaration's name does not contain the original type's name, the renaming will fail, leading to duplicate declarations.
  This will be fixed using name mangling.
- Only the type's name is used in renaming, not the package name.
  This will lead to duplicate declarations if two concrete types with the same name but
  from two different packages are used.
- Rei does not validate that the concrete type satisfies the generic type. If it does not, you will get a compilation error.
