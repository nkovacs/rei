package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// parseMapping parses a type mapping string into
// a map of generic placeholder type names and concrete Types.
// A type mapping string consists of comma separated values in the
// form Type=ConcreteType.
func parseMapping(s string) (map[string]*Type, error) {
	ret := make(map[string]*Type)
	mappings := strings.Split(s, ",")
	for _, mapping := range mappings {
		parts := strings.SplitN(mapping, "=", 2)
		if len(parts) != 2 {
			return ret, fmt.Errorf("invalid mapping %v, expected Type=ConcreteType", mapping)
		}
		if ok, _ := isIdentifier(parts[0]); !ok {
			return ret, fmt.Errorf("invalid mapping %v, %v is not a valid identifier", mapping, parts[0])
		}
		tp, err := ParseType(parts[1])
		if err != nil {
			return ret, errors.Wrap(err, fmt.Sprintf("in mapping %v", mapping))
		}
		if _, ok := ret[parts[0]]; ok {
			// This is currently not supported
			return ret, fmt.Errorf("duplicate mapping for template type %v", parts[0])
		}
		ret[parts[0]] = &tp
	}
	return ret, nil
}
