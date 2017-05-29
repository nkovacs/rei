package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Type describes a concrete type.
type Type struct {
	Pkg     string
	PkgName string
	Name    string
}

func isIdentifier(s string) (bool, int) {
	if len(s) == 0 {
		return false, -1
	}
	firstRune := []rune(s)[0]
	if !unicode.IsLetter(firstRune) && firstRune != '_' {
		return false, 0
	}
	firstError := strings.IndexFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_'
	})
	if firstError != -1 {
		return false, firstError
	}
	return true, -1
}

func validateType(t Type) (Type, error) {
	if ok, idx := isIdentifier(t.Name); !ok {
		return t, fmt.Errorf("invalid type: %v (at %v)", t.Name, idx)
	}
	if ok, idx := isIdentifier(t.PkgName); !ok && idx != -1 {
		return t, fmt.Errorf("invalid package name: %v (at %v)", t.PkgName, idx)
	}
	return t, nil
}

// ParseType parses a type string.
// The following formats are accepted:
// ConcreteType
// pkg/pkg/pkg.ConcreteType
// ("pkg/pkg/go-pkg")pkg.ConcreteType
func ParseType(s string) (Type, error) {
	dotIdx := strings.LastIndex(s, ".")
	if dotIdx == -1 {
		return validateType(Type{
			Pkg:     "",
			PkgName: "",
			Name:    s,
		})
	}
	packagePart := s[:dotIdx]
	typeName := s[dotIdx+1:]
	pkgImport := packagePart
	pkgName := packagePart

	if strings.HasPrefix(packagePart, `("`) {
		closeIdx := strings.LastIndex(packagePart, `")`)
		if closeIdx == -1 {
			return Type{
				Pkg:     packagePart,
				PkgName: packagePart,
				Name:    typeName,
			}, fmt.Errorf(`invalid type specification %v: missing closing ")`, s)
		}
		pkgImport = packagePart[2:closeIdx]
		pkgName = packagePart[closeIdx+2:]
	} else if slashIdx := strings.LastIndex(packagePart, "/"); slashIdx != -1 {
		pkgName = packagePart[slashIdx+1:]
	}

	return validateType(Type{
		Pkg:     pkgImport,
		PkgName: pkgName,
		Name:    typeName,
	})
}
