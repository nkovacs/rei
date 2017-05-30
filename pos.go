package main

import (
	"fmt"
	"go/ast"
	"go/token"
)

type tokenPositioner struct {
	currentPos token.Pos
}

func newTokenPositioner() *tokenPositioner {
	return &tokenPositioner{
		currentPos: 1,
	}
}

func (t *tokenPositioner) next() token.Pos {
	return t.nextN(1)
}

func (t *tokenPositioner) nextN(n int) token.Pos {
	ret := t.currentPos
	t.currentPos += token.Pos(n)
	return ret
}

func (t *tokenPositioner) fixPositions(n ast.Node) {
	if n == nil {
		return
	}
	switch n := n.(type) {
	case *ast.Comment:
		if n == nil {
			return
		}
		n.Slash = t.nextN(len(n.Text))
	case *ast.CommentGroup:
		if n == nil {
			return
		}
		for _, c := range n.List {
			t.fixPositions(c)
		}
	case *ast.Field:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		for _, name := range n.Names {
			t.fixPositions(name)
			t.nextN(1) // ,
		}
		t.fixPositions(n.Type)
		t.fixPositions(n.Tag)
		t.fixPositions(n.Comment)
	case *ast.FieldList:
		if n == nil {
			return
		}
		if len(n.List) == 0 {
			n.Opening = token.NoPos
			n.Closing = token.NoPos
		} else {
			n.Opening = t.next()
			for _, f := range n.List {
				t.fixPositions(f)
				t.nextN(1) // ,
			}
			n.Closing = t.next()
		}
	case *ast.BadExpr:
		if n == nil {
			return
		}
		n.From = t.nextN(int(n.To - n.From))
		n.To = t.next()
	case *ast.Ident:
		if n == nil {
			return
		}
		n.NamePos = t.nextN(len(n.Name))
	case *ast.Ellipsis:
		if n == nil {
			return
		}
		if n.Elt != nil {
			n.Ellipsis = t.next()
			t.fixPositions(n.Elt)
		} else {
			n.Ellipsis = t.nextN(3)
		}
	case *ast.BasicLit:
		if n == nil {
			return
		}
		n.ValuePos = t.nextN(len(n.Value))
	case *ast.FuncLit:
		if n == nil {
			return
		}
		t.fixPositions(n.Type)
		t.fixPositions(n.Body)
	case *ast.CompositeLit:
		if n == nil {
			return
		}
		if n.Type != nil {
			t.fixPositions(n.Type)
		}
		n.Lbrace = t.next()
		for _, el := range n.Elts {
			t.fixPositions(el)
			t.nextN(1) // ,
		}
		n.Rbrace = t.next()
	case *ast.ParenExpr:
		if n == nil {
			return
		}
		n.Lparen = t.next()
		t.fixPositions(n.X)
		n.Rparen = t.next()
	case *ast.SelectorExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		t.fixPositions(n.Sel)
	case *ast.IndexExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		n.Lbrack = t.next()
		t.fixPositions(n.Index)
		n.Rbrack = t.next()
	case *ast.SliceExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		n.Lbrack = t.next()
		t.fixPositions(n.Low)
		t.fixPositions(n.High)
		t.fixPositions(n.Max)
		n.Rbrack = t.next()
	case *ast.TypeAssertExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		n.Lparen = t.next()
		if n.Type != nil {
			t.fixPositions(n.Type)
		} else {
			t.nextN(4)
		}
		n.Rparen = t.next()
	case *ast.CallExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.Fun)
		n.Lparen = t.next()
		for _, arg := range n.Args {
			t.fixPositions(arg)
			t.nextN(1) // ,
		}
		if n.Ellipsis.IsValid() {
			n.Ellipsis = t.nextN(3)
		}
		n.Rparen = t.next()
	case *ast.StarExpr:
		if n == nil {
			return
		}
		n.Star = t.next()
		t.fixPositions(n.X)
	case *ast.UnaryExpr:
		if n == nil {
			return
		}
		n.OpPos = t.nextN(len(n.Op.String()))
		t.fixPositions(n.X)
	case *ast.BinaryExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		n.OpPos = t.nextN(len(n.Op.String()))
		t.fixPositions(n.Y)
	case *ast.KeyValueExpr:
		if n == nil {
			return
		}
		t.fixPositions(n.Key)
		n.Colon = t.next()
		t.fixPositions(n.Value)
	case *ast.ArrayType:
		if n == nil {
			return
		}
		n.Lbrack = t.next()
		t.fixPositions(n.Len)
		t.next() // closing bracket
		t.fixPositions(n.Elt)
	case *ast.StructType:
		if n == nil {
			return
		}
		n.Struct = t.nextN(len("struct"))
		t.fixPositions(n.Fields)
	case *ast.FuncType:
		if n == nil {
			return
		}
		if n.Func.IsValid() {
			n.Func = t.nextN(len("func"))
		}
		t.fixPositions(n.Params)
		t.next()
		t.fixPositions(n.Results)
	case *ast.InterfaceType:
		if n == nil {
			return
		}
		n.Interface = t.nextN(len("interface"))
		t.fixPositions(n.Methods)
	case *ast.MapType:
		if n == nil {
			return
		}
		n.Map = t.nextN(len("map"))
		t.next() // [
		t.fixPositions(n.Key)
		t.next() // ]
		t.fixPositions(n.Value)
	case *ast.ChanType:
		if n == nil {
			return
		}
		if n.Dir == ast.RECV {
			// <-chan type
			n.Begin = t.nextN(len("<-chan "))
			n.Arrow = n.Begin
			t.fixPositions(n.Value)
		} else if n.Dir == ast.SEND {
			// chan<- type
			n.Begin = t.nextN(len("chan"))
			n.Arrow = t.nextN(len("<- "))
			t.fixPositions(n.Value)
		} else {
			// chan type
			n.Begin = t.nextN(len("chan "))
			n.Arrow = token.NoPos
			t.fixPositions(n.Value)
		}
	case *ast.BadStmt:
		if n == nil {
			return
		}
		n.From = t.nextN(int(n.To - n.From))
		n.To = t.next()
	case *ast.DeclStmt:
		if n == nil {
			return
		}
		t.fixPositions(n.Decl)
	case *ast.EmptyStmt:
		if n == nil {
			return
		}
		if n.Implicit {
			n.Semicolon = t.nextN(0)
		} else {
			n.Semicolon = t.nextN(1)
		}
	case *ast.LabeledStmt:
		if n == nil {
			return
		}
		t.fixPositions(n.Label)
		n.Colon = t.next()
		t.fixPositions(n.Stmt)
	case *ast.ExprStmt:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
	case *ast.SendStmt:
		if n == nil {
			return
		}
		t.fixPositions(n.Chan)
		n.Arrow = t.nextN(len("<-"))
		t.fixPositions(n.Value)
	case *ast.IncDecStmt:
		if n == nil {
			return
		}
		t.fixPositions(n.X)
		n.TokPos = t.nextN(len("++"))
	case *ast.AssignStmt:
		if n == nil {
			return
		}
		for _, x := range n.Lhs {
			t.fixPositions(x)
			t.nextN(1) // ,
		}
		n.TokPos = t.nextN(len(n.Tok.String()))
		for _, x := range n.Rhs {
			t.fixPositions(x)
			t.nextN(1) // ,
		}
	case *ast.GoStmt:
		if n == nil {
			return
		}
		n.Go = t.nextN(len("go "))
		t.fixPositions(n.Call)
	case *ast.DeferStmt:
		if n == nil {
			return
		}
		n.Defer = t.nextN(len("defer "))
		t.fixPositions(n.Call)
	case *ast.ReturnStmt:
		if n == nil {
			return
		}
		n.Return = t.nextN(len("return "))
		for _, x := range n.Results {
			t.fixPositions(x)
			t.nextN(1) // ,
		}
	case *ast.BranchStmt:
		if n == nil {
			return
		}
		n.TokPos = t.nextN(len(n.Tok.String()))
		if n.Label != nil {
			t.next() // space
			t.fixPositions(n.Label)
		}
	case *ast.BlockStmt:
		if n == nil {
			return
		}
		n.Lbrace = t.next()
		for _, stmt := range n.List {
			t.fixPositions(stmt)
			t.nextN(1)
		}
		n.Rbrace = t.next()
	case *ast.IfStmt:
		if n == nil {
			return
		}
		n.If = t.nextN(len("if "))
		if n.Init != nil {
			t.fixPositions(n.Init)
			t.nextN(1) // ;
		}
		t.fixPositions(n.Cond)
		t.fixPositions(n.Body)
		if n.Else != nil {
			t.nextN(len(" else "))
			t.fixPositions(n.Else)
		}
	case *ast.CaseClause:
		if n == nil {
			return
		}
		if len(n.List) == 0 {
			n.Case = t.nextN(len("default"))
		} else {
			n.Case = t.nextN(len("case "))
			for _, x := range n.List {
				t.fixPositions(x)
				t.nextN(1) // ,
			}
		}
		n.Colon = t.next()
		for _, stmt := range n.Body {
			t.fixPositions(stmt)
		}
	case *ast.SwitchStmt:
		if n == nil {
			return
		}
		n.Switch = t.nextN(len("switch "))
		if n.Init != nil {
			t.fixPositions(n.Init)
			t.nextN(1) // ;
		}
		if n.Tag != nil {
			t.fixPositions(n.Tag)
			t.nextN(1) // ;
		}
		t.fixPositions(n.Body)
	case *ast.TypeSwitchStmt:
		if n == nil {
			return
		}
		n.Switch = t.nextN(len("switch "))
		if n.Init != nil {
			t.fixPositions(n.Init)
			t.nextN(1) // ;
		}
		if n.Assign != nil {
			t.fixPositions(n.Assign)
			t.nextN(1) // ;
		}
		t.fixPositions(n.Body)
	case *ast.CommClause:
		if n == nil {
			return
		}
		if n.Comm == nil {
			n.Case = t.nextN(len("default"))
		} else {
			n.Case = t.nextN(len("case "))
			t.fixPositions(n.Comm)
		}
		n.Colon = t.next()
		for _, stmt := range n.Body {
			t.fixPositions(stmt)
		}
	case *ast.SelectStmt:
		if n == nil {
			return
		}
		n.Select = t.nextN(len("select "))
		t.fixPositions(n.Body)
	case *ast.ForStmt:
		if n == nil {
			return
		}
		n.For = t.nextN(len("for "))
		if n.Init != nil {
			t.fixPositions(n.Init)
			t.nextN(1) // ;
		}
		if n.Cond != nil {
			t.fixPositions(n.Cond)
			t.nextN(1) // ;
		}
		if n.Post != nil {
			t.fixPositions(n.Post)
			t.nextN(1) // ;
		}
		t.fixPositions(n.Body)
	case *ast.RangeStmt:
		if n == nil {
			return
		}
		n.For = t.nextN(len("for "))
		t.fixPositions(n.Key)
		t.nextN(1) // ,
		t.fixPositions(n.Value)
		if n.Key != nil {
			n.TokPos = t.nextN(len(n.Tok.String()))
		}
		t.nextN(len(" range "))
		t.fixPositions(n.X)
		t.fixPositions(n.Body)
	case *ast.ImportSpec:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		t.fixPositions(n.Name)
		t.nextN(1)
		t.fixPositions(n.Path)
		t.fixPositions(n.Comment)
		n.EndPos = token.NoPos
	case *ast.ValueSpec:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		for _, name := range n.Names {
			t.fixPositions(name)
			t.nextN(1) // ,
		}
		t.fixPositions(n.Type)
		if len(n.Values) > 0 {
			t.nextN(1) // =
			for _, val := range n.Values {
				t.fixPositions(val)
				t.nextN(1) // ,
			}
		}
		t.fixPositions(n.Comment)
	case *ast.TypeSpec:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		t.fixPositions(n.Name)
		t.nextN(1)
		t.fixPositions(n.Type)
		t.fixPositions(n.Comment)
	case *ast.BadDecl:
		if n == nil {
			return
		}
		n.From = t.nextN(int(n.To - n.From))
		n.To = t.next()
	case *ast.GenDecl:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		n.TokPos = t.nextN(len(n.Tok.String()))
		t.next() // space
		if n.Lparen.IsValid() {
			n.Lparen = t.next()
		}
		for _, spec := range n.Specs {
			t.fixPositions(spec)
			t.nextN(1) // just in case
		}
		if n.Lparen.IsValid() {
			n.Rparen = t.next()
		}
	case *ast.FuncDecl:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		n.Type.Func = t.nextN(len("func "))
		t.fixPositions(n.Recv)
		t.next()
		t.fixPositions(n.Name)
		t.fixPositions(n.Type.Params)
		t.next()
		t.fixPositions(n.Type.Results)
		t.fixPositions(n.Body)
	case *ast.Package:
		if n == nil {
			return
		}
		// nothing to do
	case *ast.File:
		if n == nil {
			return
		}
		t.fixPositions(n.Doc)
		n.Package = t.nextN(len("package "))
		t.fixPositions(n.Name)
		for _, decl := range n.Decls {
			t.fixPositions(decl)
			t.next()
		}
	default:
		panic(fmt.Sprintf("unknown node: %#v", n))
	}
}

func clearPositions(n ast.Node) {
	if n == nil {
		return
	}
	switch n := n.(type) {
	case *ast.Comment:
		if n == nil {
			return
		}
		n.Slash = 0
	case *ast.CommentGroup:
		if n == nil {
			return
		}
		for _, c := range n.List {
			clearPositions(c)
		}
	case *ast.Field:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		for _, name := range n.Names {
			clearPositions(name)
		}
		clearPositions(n.Type)
		clearPositions(n.Tag)
		clearPositions(n.Comment)
	case *ast.FieldList:
		if n == nil {
			return
		}
		if n.Opening.IsValid() {
			n.Opening = 0
		}
		if n.Closing.IsValid() {
			n.Closing = 0
		}
		for _, f := range n.List {
			clearPositions(f)
		}
	case *ast.BadExpr:
		if n == nil {
			return
		}
		n.From = 0
		n.To = 0
	case *ast.Ident:
		if n == nil {
			return
		}
		n.NamePos = 0
	case *ast.Ellipsis:
		if n == nil {
			return
		}
		n.Ellipsis = 0
	case *ast.BasicLit:
		if n == nil {
			return
		}
		n.ValuePos = 0
	case *ast.FuncLit:
		if n == nil {
			return
		}
		clearPositions(n.Type)
		clearPositions(n.Body)
	case *ast.CompositeLit:
		if n == nil {
			return
		}
		if n.Type != nil {
			clearPositions(n.Type)
		}
		n.Lbrace = 0
		for _, el := range n.Elts {
			clearPositions(el)
		}
		n.Rbrace = 0
	case *ast.ParenExpr:
		if n == nil {
			return
		}
		n.Lparen = 0
		clearPositions(n.X)
		n.Rparen = 0
	case *ast.SelectorExpr:
		if n == nil {
			return
		}
		clearPositions(n.X)
		clearPositions(n.Sel)
	case *ast.IndexExpr:
		if n == nil {
			return
		}
		clearPositions(n.X)
		n.Lbrack = 0
		clearPositions(n.Index)
		n.Rbrack = 0
	case *ast.SliceExpr:
		if n == nil {
			return
		}
		clearPositions(n.X)
		n.Lbrack = 0
		clearPositions(n.Low)
		clearPositions(n.High)
		clearPositions(n.Max)
		n.Rbrack = 0
	case *ast.TypeAssertExpr:
		if n == nil {
			return
		}
		clearPositions(n.X)
		n.Lparen = 0
		if n.Type != nil {
			clearPositions(n.Type)
		}
		n.Rparen = 0
	case *ast.CallExpr:
		if n == nil {
			return
		}
		clearPositions(n.Fun)
		n.Lparen = 0
		for _, arg := range n.Args {
			clearPositions(arg)
		}
		if n.Ellipsis.IsValid() {
			n.Ellipsis = 1
		}
		n.Rparen = 0
	case *ast.StarExpr:
		if n == nil {
			return
		}
		n.Star = 0
		clearPositions(n.X)
	case *ast.UnaryExpr:
		if n == nil {
			return
		}
		n.OpPos = 0
		clearPositions(n.X)
	case *ast.BinaryExpr:
		if n == nil {
			return
		}
		clearPositions(n.X)
		n.OpPos = 0
		clearPositions(n.Y)
	case *ast.KeyValueExpr:
		if n == nil {
			return
		}
		clearPositions(n.Key)
		n.Colon = 0
		clearPositions(n.Value)
	case *ast.ArrayType:
		if n == nil {
			return
		}
		n.Lbrack = 0
		clearPositions(n.Len)
		clearPositions(n.Elt)
	case *ast.StructType:
		if n == nil {
			return
		}
		n.Struct = 0
		clearPositions(n.Fields)
	case *ast.FuncType:
		if n == nil {
			return
		}
		if n.Func.IsValid() {
			n.Func = 1
		}
		clearPositions(n.Params)
		clearPositions(n.Results)
	case *ast.InterfaceType:
		if n == nil {
			return
		}
		n.Interface = 0
		clearPositions(n.Methods)
	case *ast.MapType:
		if n == nil {
			return
		}
		n.Map = 0
		clearPositions(n.Key)
		clearPositions(n.Value)
	case *ast.ChanType:
		if n == nil {
			return
		}
		n.Begin = 0
		if n.Arrow.IsValid() {
			n.Arrow = 1
		}
		clearPositions(n.Value)
	case *ast.BadStmt:
		if n == nil {
			return
		}
		n.From = 0
		n.To = 0
	case *ast.DeclStmt:
		if n == nil {
			return
		}
		clearPositions(n.Decl)
	case *ast.EmptyStmt:
		if n == nil {
			return
		}
		n.Semicolon = 0
	case *ast.LabeledStmt:
		if n == nil {
			return
		}
		clearPositions(n.Label)
		n.Colon = 0
		clearPositions(n.Stmt)
	case *ast.ExprStmt:
		if n == nil {
			return
		}
		clearPositions(n.X)
	case *ast.SendStmt:
		if n == nil {
			return
		}
		clearPositions(n.Chan)
		n.Arrow = 0
		clearPositions(n.Value)
	case *ast.IncDecStmt:
		if n == nil {
			return
		}
		clearPositions(n.X)
		n.TokPos = 0
	case *ast.AssignStmt:
		if n == nil {
			return
		}
		for _, x := range n.Lhs {
			clearPositions(x)
		}
		n.TokPos = 0
		for _, x := range n.Rhs {
			clearPositions(x)
		}
	case *ast.GoStmt:
		if n == nil {
			return
		}
		n.Go = 0
		clearPositions(n.Call)
	case *ast.DeferStmt:
		if n == nil {
			return
		}
		n.Defer = 0
		clearPositions(n.Call)
	case *ast.ReturnStmt:
		if n == nil {
			return
		}
		n.Return = 0
		for _, x := range n.Results {
			clearPositions(x)
		}
	case *ast.BranchStmt:
		if n == nil {
			return
		}
		n.TokPos = 0
		if n.Label != nil {
			clearPositions(n.Label)
		}
	case *ast.BlockStmt:
		if n == nil {
			return
		}
		n.Lbrace = 0
		for _, stmt := range n.List {
			clearPositions(stmt)
		}
		n.Rbrace = 0
	case *ast.IfStmt:
		if n == nil {
			return
		}
		n.If = 0
		if n.Init != nil {
			clearPositions(n.Init)
		}
		clearPositions(n.Cond)
		clearPositions(n.Body)
		if n.Else != nil {
			clearPositions(n.Else)
		}
	case *ast.CaseClause:
		if n == nil {
			return
		}
		n.Case = 0
		for _, x := range n.List {
			clearPositions(x)
		}
		n.Colon = 0
		for _, stmt := range n.Body {
			clearPositions(stmt)
		}
	case *ast.SwitchStmt:
		if n == nil {
			return
		}
		n.Switch = 0
		if n.Init != nil {
			clearPositions(n.Init)
		}
		if n.Tag != nil {
			clearPositions(n.Tag)
		}
		clearPositions(n.Body)
	case *ast.TypeSwitchStmt:
		if n == nil {
			return
		}
		n.Switch = 0
		if n.Init != nil {
			clearPositions(n.Init)
		}
		if n.Assign != nil {
			clearPositions(n.Assign)
		}
		clearPositions(n.Body)
	case *ast.CommClause:
		if n == nil {
			return
		}
		n.Case = 0
		clearPositions(n.Comm)
		n.Colon = 0
		for _, stmt := range n.Body {
			clearPositions(stmt)
		}
	case *ast.SelectStmt:
		if n == nil {
			return
		}
		n.Select = 0
		clearPositions(n.Body)
	case *ast.ForStmt:
		if n == nil {
			return
		}
		n.For = 0
		if n.Init != nil {
			clearPositions(n.Init)
		}
		if n.Cond != nil {
			clearPositions(n.Cond)
		}
		if n.Post != nil {
			clearPositions(n.Post)
		}
		clearPositions(n.Body)
	case *ast.RangeStmt:
		if n == nil {
			return
		}
		n.For = 0
		clearPositions(n.Key)
		clearPositions(n.Value)
		if n.Key != nil {
			n.TokPos = 0
		}
		clearPositions(n.X)
		clearPositions(n.Body)
	case *ast.ImportSpec:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		clearPositions(n.Name)
		clearPositions(n.Path)
		clearPositions(n.Comment)
		n.EndPos = token.NoPos
	case *ast.ValueSpec:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		for _, name := range n.Names {
			clearPositions(name)
		}
		clearPositions(n.Type)
		if len(n.Values) > 0 {
			for _, val := range n.Values {
				clearPositions(val)
			}
		}
		clearPositions(n.Comment)
	case *ast.TypeSpec:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		clearPositions(n.Name)
		clearPositions(n.Type)
		clearPositions(n.Comment)
	case *ast.BadDecl:
		if n == nil {
			return
		}
		n.From = 0
		n.To = 0
	case *ast.GenDecl:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		n.TokPos = 0
		if n.Lparen.IsValid() {
			n.Lparen = 0
		}
		for _, spec := range n.Specs {
			clearPositions(spec)
		}
		if n.Lparen.IsValid() {
			n.Rparen = 0
		}
		if len(n.Specs) > 1 {
			n.Lparen = 1
		}
	case *ast.FuncDecl:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		n.Type.Func = 0
		clearPositions(n.Recv)
		clearPositions(n.Name)
		clearPositions(n.Type.Params)
		clearPositions(n.Type.Results)
		clearPositions(n.Body)
	case *ast.Package:
		if n == nil {
			return
		}
		// nothing to do
	case *ast.File:
		if n == nil {
			return
		}
		clearPositions(n.Doc)
		n.Package = 0
		clearPositions(n.Name)
		for _, decl := range n.Decls {
			clearPositions(decl)
		}
	default:
		panic(fmt.Sprintf("unknown node: %#v", n))
	}
}
