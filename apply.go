// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"reflect"
)

// An ApplyFunc is invoked by Apply for each node n, even if n is nil,
// before and/or after the node's children.
//
// The parent, name, and index arguments identify the parent node's field
// containing n. If that field is a slice, index identifies the node's position
// in that slice; index is < 0 otherwise. Roughly speaking, the following
// invariants hold:
//
//   parent.name        == n  if index < 0
//   parent.name[index] == n  if index >= 0
//
// SetField(parent, name, index, n1) can be used to change that field
// to a different node n1.
//
// Exception: If the parent is a *Package, and Apply is iterating
// through the Files map, name is the filename, and index is -1.
//
// The return value of ApplyFunc controls the syntax tree traversal.
// See Apply for details.
type ApplyFunc func(parent ast.Node, name string, index int, n ast.Node) bool

// Apply traverses a syntax tree recursively, starting with root,
// and calling pre and post for each node as described below. The
// result is the (possibly modified) syntax tree.
//
// If pre is not nil, it is called for each node before its children
// are traversed (pre-order). If the result of calling pre is false,
// no children are traversed, and post is not called for that node.
//
// If post is not nil, it is called for each node after its children
// were traversed (post-order). If the result of calling post is false,
// traversal is terminated and Apply returns immediately.
//
// Only fields that refer to AST nodes are considered children.
// Children are traversed in the order in which they appear in the
// respective node's struct definition.
func Apply(root ast.Node, pre, post ApplyFunc) ast.Node {
	defer func() {
		if r := recover(); r != nil && r != abort {
			panic(r)
		}
	}()
	a := &application{root, pre, post}
	a.apply(a, "Node", -1, a.Node)
	return a.Node
}

// SetField sets the named field in the parent node to n. If the field
// is a slice, index is the slice index. The named field must exist in
// the parent, n must be assignable to that field, and the field must be
// indexable if index >= 0. In other words, SetField performs the following
// assignment:
//
//   parent.name        = n  if index < 0
//   parent.name[index] = n  if index >= 0
//
// The parent node may be a pointer to the struct containing the named
// field, or it may be the struct itself.
//
// Exception: If the parent is a Package, n must be a *File and name is
// interpreted as the filename in the Package.Files map.
func SetField(parent ast.Node, name string, index int, n ast.Node) {
	// TODO(gri) This doesn't handle the Package.Files map yet.
	v := reflect.Indirect(reflect.ValueOf(parent)).FieldByName(name)
	if index >= 0 {
		v = v.Index(index)
	}
	v.Set(reflect.ValueOf(n))
}

type application struct {
	ast.Node
	pre, post ApplyFunc
}

func (a *application) apply(parent ast.Node, name string, index int, n ast.Node) {
	if a.pre != nil && !a.pre(parent, name, index, n) {
		return
	}

	// walk children
	// (the order of the cases matches the order
	// of the corresponding node types in ast.go)
	switch n := n.(type) {
	case nil:
		// nothing to do

	// Comments and fields
	case *ast.Comment:
		// nothing to do

	case *ast.CommentGroup:
		if n != nil {
			for i, x := range n.List {
				a.apply(n, "List", i, x)
			}
		}

	case *ast.Field:
		a.apply(n, "Doc", -1, n.Doc)
		a.applyIdentList(n, "Names", n.Names)
		a.apply(n, "Type", -1, n.Type)
		a.apply(n, "Tag", -1, n.Tag)
		a.apply(n, "Comment", -1, n.Comment)

	case *ast.FieldList:
		if n != nil {
			for i, x := range n.List {
				a.apply(n, "List", i, x)
			}
		}

	// Expressions
	case *ast.BadExpr, *ast.Ident, *ast.BasicLit:
		// nothing to do

	case *ast.Ellipsis:
		a.apply(n, "Elt", -1, n.Elt)

	case *ast.FuncLit:
		a.apply(n, "Type", -1, n.Type)
		a.apply(n, "Body", -1, n.Body)

	case *ast.CompositeLit:
		a.apply(n, "Type", -1, n.Type)
		a.applyExprList(n, "Elts", n.Elts)

	case *ast.ParenExpr:
		a.apply(n, "X", -1, n.X)

	case *ast.SelectorExpr:
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Sel", -1, n.Sel)

	case *ast.IndexExpr:
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Index", -1, n.Index)

	case *ast.SliceExpr:
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Low", -1, n.Low)
		a.apply(n, "High", -1, n.High)
		a.apply(n, "Max", -1, n.Max)

	case *ast.TypeAssertExpr:
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Type", -1, n.Type)

	case *ast.CallExpr:
		a.apply(n, "Fun", -1, n.Fun)
		a.applyExprList(n, "Args", n.Args)

	case *ast.StarExpr:
		a.apply(n, "X", -1, n.X)

	case *ast.UnaryExpr:
		a.apply(n, "X", -1, n.X)

	case *ast.BinaryExpr:
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Y", -1, n.Y)

	case *ast.KeyValueExpr:
		a.apply(n, "Key", -1, n.Key)
		a.apply(n, "Value", -1, n.Value)

	// Types
	case *ast.ArrayType:
		a.apply(n, "Len", -1, n.Len)
		a.apply(n, "Elt", -1, n.Elt)

	case *ast.StructType:
		a.apply(n, "Fields", -1, n.Fields)

	case *ast.FuncType:
		a.apply(n, "Params", -1, n.Params)
		a.apply(n, "Results", -1, n.Results)

	case *ast.InterfaceType:
		a.apply(n, "Methods", -1, n.Methods)

	case *ast.MapType:
		a.apply(n, "Key", -1, n.Key)
		a.apply(n, "Value", -1, n.Value)

	case *ast.ChanType:
		a.apply(n, "Value", -1, n.Value)

	// Statements
	case *ast.BadStmt:
		// nothing to do

	case *ast.DeclStmt:
		a.apply(n, "Decl", -1, n.Decl)

	case *ast.EmptyStmt:
		// nothing to do

	case *ast.LabeledStmt:
		a.apply(n, "Label", -1, n.Label)
		a.apply(n, "Stmt", -1, n.Stmt)

	case *ast.ExprStmt:
		a.apply(n, "X", -1, n.X)

	case *ast.SendStmt:
		a.apply(n, "Chan", -1, n.Chan)
		a.apply(n, "Value", -1, n.Value)

	case *ast.IncDecStmt:
		a.apply(n, "X", -1, n.X)

	case *ast.AssignStmt:
		a.applyExprList(n, "Lhs", n.Lhs)
		a.applyExprList(n, "Rhs", n.Rhs)

	case *ast.GoStmt:
		a.apply(n, "Call", -1, n.Call)

	case *ast.DeferStmt:
		a.apply(n, "Call", -1, n.Call)

	case *ast.ReturnStmt:
		a.applyExprList(n, "Results", n.Results)

	case *ast.BranchStmt:
		a.apply(n, "Label", -1, n.Label)

	case *ast.BlockStmt:
		a.applyStmtList(n, "List", n.List)

	case *ast.IfStmt:
		a.apply(n, "Init", -1, n.Init)
		a.apply(n, "Cond", -1, n.Cond)
		a.apply(n, "Body", -1, n.Body)
		a.apply(n, "Else", -1, n.Else)

	case *ast.CaseClause:
		a.applyExprList(n, "List", n.List)
		a.applyStmtList(n, "Body", n.Body)

	case *ast.SwitchStmt:
		a.apply(n, "Init", -1, n.Init)
		a.apply(n, "Tag", -1, n.Tag)
		a.apply(n, "Body", -1, n.Body)

	case *ast.TypeSwitchStmt:
		a.apply(n, "Init", -1, n.Init)
		a.apply(n, "Assign", -1, n.Assign)
		a.apply(n, "Body", -1, n.Body)

	case *ast.CommClause:
		a.apply(n, "Comm", -1, n.Comm)
		a.applyStmtList(n, "Body", n.Body)

	case *ast.SelectStmt:
		a.apply(n, "Body", -1, n.Body)

	case *ast.ForStmt:
		a.apply(n, "Init", -1, n.Init)
		a.apply(n, "Cond", -1, n.Cond)
		a.apply(n, "Post", -1, n.Post)
		a.apply(n, "Body", -1, n.Body)

	case *ast.RangeStmt:
		a.apply(n, "Key", -1, n.Key)
		a.apply(n, "Value", -1, n.Value)
		a.apply(n, "X", -1, n.X)
		a.apply(n, "Body", -1, n.Body)

	// Declarations
	case *ast.ImportSpec:
		a.apply(n, "Doc", -1, n.Doc)
		a.apply(n, "Name", -1, n.Name)
		a.apply(n, "Path", -1, n.Path)
		a.apply(n, "Comment", -1, n.Comment)

	case *ast.ValueSpec:
		a.apply(n, "Doc", -1, n.Doc)
		a.applyIdentList(n, "Names", n.Names)
		a.apply(n, "Type", -1, n.Type)
		a.applyExprList(n, "Values", n.Values)
		a.apply(n, "Comment", -1, n.Comment)

	case *ast.TypeSpec:
		a.apply(n, "Doc", -1, n.Doc)
		a.apply(n, "Name", -1, n.Name)
		a.apply(n, "Type", -1, n.Type)
		a.apply(n, "Comment", -1, n.Comment)

	case *ast.BadDecl:
		// nothing to do

	case *ast.GenDecl:
		a.apply(n, "Doc", -1, n.Doc)
		for i, x := range n.Specs {
			a.apply(n, "Specs", i, x)
		}

	case *ast.FuncDecl:
		a.apply(n, "Doc", -1, n.Doc)
		a.apply(n, "Recv", -1, n.Recv)
		a.apply(n, "Name", -1, n.Name)
		a.apply(n, "Type", -1, n.Type)
		a.apply(n, "Body", -1, n.Body)

	// Files and packages
	case *ast.File:
		a.apply(n, "Doc", -1, n.Doc)
		a.apply(n, "Name", -1, n.Name)
		a.applyDeclList(n, "Decls", n.Decls)
		// don't walk n.Comments - they have been
		// visited already through the individual
		// nodes

	case *ast.Package:
		for name, f := range n.Files {
			a.apply(n, name, -1, f)
		}

	default:
		panic(fmt.Sprintf("ast.Apply: unexpected node type %T", n))
	}

	if a.post != nil && !a.post(parent, name, index, n) {
		panic(abort)
	}
}

var abort = new(int) // singleton, to signal abortion of Apply

func (a *application) applyIdentList(parent ast.Node, name string, list []*ast.Ident) {
	for i, x := range list {
		a.apply(parent, name, i, x)
	}
}

func (a *application) applyExprList(parent ast.Node, name string, list []ast.Expr) {
	for i, x := range list {
		a.apply(parent, name, i, x)
	}
}

func (a *application) applyStmtList(parent ast.Node, name string, list []ast.Stmt) {
	for i, x := range list {
		a.apply(parent, name, i, x)
	}
}

func (a *application) applyDeclList(parent ast.Node, name string, list []ast.Decl) {
	for i, x := range list {
		a.apply(parent, name, i, x)
	}
}
