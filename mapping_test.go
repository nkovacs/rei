package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseMapping(t *testing.T) {
	assert := assert.New(t)
	testCases := []struct {
		input    string
		ok       bool
		expected map[string]*Type
	}{
		{
			"Type=Concrete",
			true,
			map[string]*Type{
				"Type": {
					Pkg:     "",
					PkgName: "",
					Name:    "Concrete",
				},
			},
		},
		{
			"Type=Concrete1,Type=Concrete2",
			false,
			map[string]*Type{
				"Type": {
					Pkg:     "",
					PkgName: "",
					Name:    "Concrete1",
				},
			},
		},
		{
			"Type->Concrete",
			false,
			map[string]*Type{
				"Type": {
					Pkg:     "",
					PkgName: "",
					Name:    "Concrete",
				},
			},
		},
		{
			"1Type=Concrete",
			false,
			map[string]*Type{
				"1Type": {
					Pkg:     "",
					PkgName: "",
					Name:    "Concrete",
				},
			},
		},
		{
			"Type=github.com/user/pkg.Concrete",
			true,
			map[string]*Type{
				"Type": {
					Pkg:     "github.com/user/pkg",
					PkgName: "pkg",
					Name:    "Concrete",
				},
			},
		},
		{
			`Type1=github.com/user/pkg.Concrete,Type2=foo/bar.Concrete,Type3=("github.com/user/go-pkg")pkg.Concrete`,
			true,
			map[string]*Type{
				"Type1": {
					Pkg:     "github.com/user/pkg",
					PkgName: "pkg",
					Name:    "Concrete",
				},
				"Type2": {
					Pkg:     "foo/bar",
					PkgName: "bar",
					Name:    "Concrete",
				},
				"Type3": {
					Pkg:     "github.com/user/go-pkg",
					PkgName: "pkg",
					Name:    "Concrete",
				},
			},
		},
		{
			"Type=github.com/user/pkgConcrete",
			false,
			map[string]*Type{
				"Type": {
					Pkg:     "github.com/user/pkg",
					PkgName: "pkg",
					Name:    "Concrete",
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			mapping, err := parseMapping(tc.input)
			assert.Equal(tc.ok, err == nil, fmt.Sprintf("%v: %v", tc.input, err))
			if tc.ok {
				assert.Equal(tc.expected, mapping, tc.input)
			}
		})
	}
}
