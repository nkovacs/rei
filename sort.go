package main

import (
	"go/ast"
	"go/token"
	"sort"
)

type SortableSpec struct {
	s   ast.Spec
	pos token.Pos
}

type SortableSpecs []*SortableSpec

type SortableDecl struct {
	d   ast.Decl
	pos token.Pos
}

type SortableDecls []*SortableDecl

// Len is the number of elements in the collection.
func (s SortableSpecs) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SortableSpecs) Less(i, j int) bool {
	return s[i].pos < s[j].pos
}

// Swap swaps the elements with indexes i and j.
func (s SortableSpecs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Len is the number of elements in the collection.
func (s SortableDecls) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s SortableDecls) Less(i, j int) bool {
	return s[i].pos < s[j].pos
}

// Swap swaps the elements with indexes i and j.
func (s SortableDecls) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortSpecs(in map[token.Pos]ast.Spec) []ast.Spec {
	sIn := make(SortableSpecs, 0, len(in))
	for p, s := range in {
		sIn = append(sIn, &SortableSpec{
			s:   s,
			pos: p,
		})
	}
	sort.Sort(sIn)
	ret := make([]ast.Spec, len(sIn))
	for i, s := range sIn {
		ret[i] = s.s
	}
	return ret
}

func sortDecls(in map[token.Pos]ast.Decl) []ast.Decl {
	sIn := make(SortableDecls, 0, len(in))
	for p, d := range in {
		sIn = append(sIn, &SortableDecl{
			d:   d,
			pos: p,
		})
	}
	sort.Sort(sIn)
	ret := make([]ast.Decl, len(sIn))
	for i, s := range sIn {
		ret[i] = s.d
	}
	return ret
}
