package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIdentifier(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		name     string
		expected bool
		errorIdx int
	}{
		{"Type", true, 0},
		{"a", true, 0},
		{"_x9", true, 0},
		{"ThisVariableIsExported", true, 0},
		{"αβ", true, 0},
		{"(foo)", false, 0},
		{"12foo", false, 0},
		{"", false, -1},
		{"foo.", false, 3},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ok, errorIdx := isIdentifier(tc.name)
			assert.Equal(tc.expected, ok, tc.name)
			if !ok {
				assert.Equal(tc.errorIdx, errorIdx, tc.name+" first error")
			}
		})
	}
}

func TestParseType(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		name     string
		ok       bool
		expected Type
	}{
		{"Concrete", true, Type{Pkg: "", PkgName: "", Name: "Concrete"}},
		{"1Concrete", false, Type{Pkg: "", PkgName: "", Name: ""}},
		{"pkg.Concrete", true, Type{Pkg: "pkg", PkgName: "pkg", Name: "Concrete"}},
		{"github.com/user/pkg/subpkg.Concrete", true, Type{Pkg: "github.com/user/pkg/subpkg", PkgName: "subpkg", Name: "Concrete"}},
		{`("github.com/user/pkg/go-subpkg")subpkg.Concrete`, true, Type{Pkg: "github.com/user/pkg/go-subpkg", PkgName: "subpkg", Name: "Concrete", Aliased: true}},

		{`(github.com/user/pkg/go-subpkg)subpkg.Concrete`, false, Type{Pkg: "github.com/user/pkg/go-subpkg", PkgName: "subpkg", Name: "Concrete", Aliased: true}},
		{`"github.com/user/pkg/go-subpkg"subpkg.Concrete`, false, Type{Pkg: "github.com/user/pkg/go-subpkg", PkgName: "subpkg", Name: "Concrete", Aliased: true}},
		{`("github.com/user/pkg/go-subpkg)subpkg.Concrete`, false, Type{Pkg: "github.com/user/pkg/go-subpkg", PkgName: "subpkg", Name: "Concrete", Aliased: true}},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			tp, err := ParseType(tc.name)
			assert.Equal(tc.ok, err == nil, fmt.Sprintf("%v: %v", tc.name, err))
			if tc.ok {
				assert.Equal(tc.expected, tp, tc.name)
			}
		})
	}
}
