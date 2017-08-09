package main

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/imports"

	"github.com/pkg/errors"
)

type empty struct{}

type genericContext struct {
	fset         *token.FileSet
	genericTypes map[string]*Type

	types     map[token.Pos]ast.Spec
	isGeneric map[token.Pos]bool
	funcs     map[token.Pos]ast.Decl
	vars      map[token.Pos]ast.Spec
	consts    map[token.Pos]ast.Spec

	renamer     *strings.Replacer
	renamePairs []string
	renames     map[token.Pos]ast.Expr //*ast.SelectorExpr or *ast.Ident or *ast.StarExpr

	visited map[token.Pos]bool
}

func lowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToLower(r)) + s[n:]
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

func (gctx *genericContext) registerGenericType(node ast.Decl) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok {
		return false
	}
	for _, spec := range decl.Specs {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		if ts.Name == nil {
			continue
		}
		gType, ok := gctx.genericTypes[ts.Name.String()]
		if !ok {
			continue
		}
		gctx.types[ts.Pos()] = ts
		gctx.isGeneric[ts.Pos()] = true
		gctx.visited[ts.Pos()] = true
		if gType.PkgName != "" {
			gctx.renames[ts.Pos()] = &ast.SelectorExpr{
				X: &ast.Ident{
					Name: gType.PkgName,
				},
				Sel: &ast.Ident{
					Name: gType.Name,
				},
			}
		} else {
			gctx.renames[ts.Pos()] = &ast.Ident{
				Name: gType.Name,
			}
		}
		if gType.Pointer {
			gctx.renames[ts.Pos()] = &ast.StarExpr{
				X: gctx.renames[ts.Pos()],
			}
		}
		gctx.renamePairs = append(gctx.renamePairs,
			lowerFirst(ts.Name.String()), lowerFirst(gType.Name),
			upperFirst(ts.Name.String()), upperFirst(gType.Name),
		)
		return true
	}
	return false
}

func (gctx *genericContext) registerGenericTypes(file *ast.File) {
	gctx.renamePairs = make([]string, 0)
	for _, decl := range file.Decls {
		gctx.registerGenericType(decl)
	}
	gctx.renamer = strings.NewReplacer(gctx.renamePairs...)
}

func (gctx *genericContext) isDependant(node ast.Node) bool {
	found := false

	// Methods on the generic types are like interfaces,
	// they should not be reified.
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		if funcDecl.Recv != nil {
			for _, field := range funcDecl.Recv.List {
				// Ident or StarExpr
				var i *ast.Ident
				switch x := field.Type.(type) {
				case *ast.StarExpr:
					i = x.X.(*ast.Ident)
				case *ast.Ident:
					i = x
				}
				if spec, ok := i.Obj.Decl.(*ast.TypeSpec); ok {
					if gctx.isGeneric[spec.Pos()] {
						return false
					}
				}
			}
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		if n, ok := n.(*ast.Ident); ok && n.Obj != nil {
			// check types
			if spec, ok := n.Obj.Decl.(*ast.TypeSpec); ok {
				for pos := range gctx.types {
					if spec.Pos() == pos {
						found = true
						return false // no need to check children
					}
				}
			}
			// check functions
			if fdecl, ok := n.Obj.Decl.(*ast.FuncDecl); ok {
				for pos := range gctx.funcs {
					if fdecl.Pos() == pos {
						found = true
						return false // no need to check children
					}
				}
			}
			// check variables
			if spec, ok := n.Obj.Decl.(*ast.ValueSpec); ok {
				for pos := range gctx.vars {
					if spec.Pos() == pos {
						found = true
						return false // no need to check children
					}
				}
				for pos := range gctx.consts {
					if spec.Pos() == pos {
						found = true
						return false // no need to check children
					}
				}
			}
		}
		return true
	})
	return found
}

func (gctx *genericContext) addDependant(n ast.Node, isConst bool) {
	switch d := n.(type) {
	case *ast.FuncDecl:
		gctx.funcs[d.Pos()] = d
		if d.Recv != nil {
			// If a method depends on a generic type,
			// the entire struct and all methods must be
			// specialized.
			for _, field := range d.Recv.List {
				// Ident or StarExpr
				var i *ast.Ident
				switch x := field.Type.(type) {
				case *ast.StarExpr:
					i = x.X.(*ast.Ident)
				case *ast.Ident:
					i = x
				}
				gctx.addDependant(i.Obj.Decl.(ast.Node), false)
			}
		} else {
			// If it's not a method, it must be renamed.
			if d.Name != nil {
				gctx.renames[d.Pos()] = &ast.Ident{
					Name: gctx.renamer.Replace(d.Name.Name),
				}
			}
		}
	case *ast.TypeSpec:
		gctx.types[d.Pos()] = d
		if d.Name != nil {
			gctx.renames[d.Pos()] = &ast.Ident{
				Name: gctx.renamer.Replace(d.Name.Name),
			}
		}
	case *ast.ValueSpec:
		if isConst {
			gctx.consts[d.Pos()] = d
		} else {
			gctx.vars[d.Pos()] = d
		}
		for _, name := range d.Names {
			if name != nil {
				gctx.renames[name.Pos()] = &ast.Ident{
					Name: gctx.renamer.Replace(name.Name),
				}
			}
		}
	default:
		panic("invalid decl")
	}
	gctx.visited[n.Pos()] = true
}

func (gctx *genericContext) collectDependants(file *ast.File) {
	changed := true

	type nodeData struct {
		n       ast.Node
		isConst bool
	}

	var nodes []nodeData
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			nodes = append(nodes, nodeData{
				n: d,
			})
		case *ast.GenDecl:
			isConst := false
			if d.Tok == token.CONST {
				isConst = true
			}
			for _, s := range d.Specs {
				switch s := s.(type) {
				case *ast.ImportSpec:
					if s.Doc == nil {
						s.Doc = d.Doc
					}
				case *ast.TypeSpec:
					if s.Doc == nil {
						s.Doc = d.Doc
					}
				case *ast.ValueSpec:
					if s.Doc == nil {
						s.Doc = d.Doc
					}
				}
				nodes = append(nodes, nodeData{
					n:       s,
					isConst: isConst,
				})
			}
		}
	}
	for changed {
		changed = false
		for _, node := range nodes {
			if gctx.visited[node.n.Pos()] {
				continue
			}
			if gctx.isDependant(node.n) {
				changed = true
				gctx.addDependant(node.n, node.isConst)
			}
		}
	}
}

func (gctx *genericContext) doRenames(n ast.Node) {
	/*
		fmt.Printf("Renames:\n")
		for k, v := range gctx.renames {
			fmt.Printf("%v: %#v\n", gctx.fset.Position(k).String(), v)
		}
		fmt.Printf("\n\n")
	*/

	// We can't touch the file until we've finished analyzing it,
	// because replacing a declaration with the renamed declaration
	// will change the original declaration's pos in the Ident's Obj.Decl field.
	type renameJob struct {
		parent      ast.Node
		name        string
		index       int
		replacement ast.Node
	}
	renames := make([]*renameJob, 0)
	Apply(n, func(parent ast.Node, name string, index int, n ast.Node) bool {
		if n, ok := n.(*ast.Ident); ok && n != nil && n.Obj != nil {
			if gctx.isGeneric[n.Pos()] {
				// Skip the generic type declaration.
				return false
			}
			// check types
			if spec, ok := n.Obj.Decl.(*ast.TypeSpec); ok {
				if renameTo, ok := gctx.renames[spec.Pos()]; ok {
					renames = append(renames, &renameJob{
						parent:      parent,
						name:        name,
						index:       index,
						replacement: renameTo,
					})
				}
			}
			// check functions
			if fdecl, ok := n.Obj.Decl.(*ast.FuncDecl); ok {
				if renameTo, ok := gctx.renames[fdecl.Pos()]; ok {
					renames = append(renames, &renameJob{
						parent:      parent,
						name:        name,
						index:       index,
						replacement: renameTo,
					})
				}
			}
			// check variables
			if spec, ok := n.Obj.Decl.(*ast.ValueSpec); ok {
				if parent.Pos() == spec.Pos() {
					if renameTo, ok := gctx.renames[n.Pos()]; ok {
						renames = append(renames, &renameJob{
							parent:      parent,
							name:        name,
							index:       index,
							replacement: renameTo,
						})
					}
				} else {
					for _, ident := range spec.Names {
						if renameTo, ok := gctx.renames[ident.Pos()]; ok && ident.Name == n.Name {
							renames = append(renames, &renameJob{
								parent:      parent,
								name:        name,
								index:       index,
								replacement: renameTo,
							})
						}
					}
				}
			}
		}
		return true
	}, nil)
	for _, job := range renames {
		SetField(job.parent, job.name, job.index, job.replacement)
	}
}

func (gctx *genericContext) renameComments(cg *ast.CommentGroup) *ast.CommentGroup {
	if cg == nil {
		return cg
	}
	for _, c := range cg.List {
		c.Text = gctx.renamer.Replace(c.Text)
	}
	return cg
}

func gen(in io.Reader, inFilename string, targetPackageName string, typeMapping map[string]*Type, out io.Writer, outFilename string) error {

	gctx := &genericContext{
		fset:         token.NewFileSet(),
		genericTypes: typeMapping,
		types:        make(map[token.Pos]ast.Spec),
		isGeneric:    make(map[token.Pos]bool),
		funcs:        make(map[token.Pos]ast.Decl),
		vars:         make(map[token.Pos]ast.Spec),
		consts:       make(map[token.Pos]ast.Spec),
		visited:      make(map[token.Pos]bool),
		renames:      make(map[token.Pos]ast.Expr),
	}
	file, err := parser.ParseFile(gctx.fset, inFilename, in, parser.ParseComments)
	if err != nil {
		return errors.Wrap(err, "parsing file failed")
	}

	// ast.Print(gctx.fset, file)

	// TODO: multiple files

	outImports := make([]*ast.ImportSpec, 0)

	for _, gType := range typeMapping {
		if gType.Pkg == "" {
			continue
		}
		found := false
		// Check if specific type's package is imported already.
		for _, importSpec := range file.Imports {
			if importSpec.Path == nil {
				continue
			}
			unquoted, err := strconv.Unquote(importSpec.Path.Value)
			if err != nil {
				return errors.Wrap(err, "could not unquote import path")
			}
			if gType.Pkg != "" && gType.Pkg == unquoted {
				found = true
				if importSpec.Name == nil {
					// Assume the package name in the typemapping is correct.
					continue
				}
				if importSpec.Name.Name == "." {
					gType.PkgName = ""
				} else {
					gType.PkgName = importSpec.Name.Name
				}
				continue
			}
		}
		if !found {
			var nameIdent *ast.Ident
			if gType.Aliased {
				nameIdent = &ast.Ident{
					Name: gType.PkgName,
				}
			}
			outImports = append(outImports, &ast.ImportSpec{
				Name: nameIdent,
				Path: &ast.BasicLit{
					Kind:  token.STRING,
					Value: strconv.Quote(gType.Pkg),
				},
			})
		}
	}

	outImports = append(outImports, file.Imports...)

	gctx.registerGenericTypes(file)
	gctx.collectDependants(file)

	/*
		fmt.Println("Dependants")
		for pos := range gctx.visited {
			fmt.Println("\t", gctx.fset.Position(pos).String())
		}
		fmt.Printf("\n\n")
	*/

	if targetPackageName == "" {
		targetPackageName = file.Name.Name
	}

	// rename types
	gctx.doRenames(file)

	// remove generic types from output
	for p := range gctx.types {
		if gctx.isGeneric[p] {
			delete(gctx.types, p)
		}
	}

	outFset := token.NewFileSet()
	outFile := &ast.File{
		Name: &ast.Ident{
			Name: targetPackageName,
		},
	}

	if len(outImports) > 0 {
		importDecl := &ast.GenDecl{
			Tok: token.IMPORT,
		}
		for _, spec := range outImports {
			importDecl.Specs = append(importDecl.Specs, spec)
		}
		outFile.Decls = append(outFile.Decls, importDecl)
	}

	sortedTypes := sortSpecs(gctx.types)
	for _, spec := range sortedTypes {
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}
		newTs := &ast.TypeSpec{
			Name: ts.Name,
			Type: ts.Type,
		}
		decl := &ast.GenDecl{
			Tok: token.TYPE,
			Doc: gctx.renameComments(ts.Doc),
			Specs: []ast.Spec{
				newTs,
			},
		}
		outFile.Decls = append(outFile.Decls, decl)
	}

	sortedConsts := sortSpecs(gctx.consts)
	for _, spec := range sortedConsts {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		newVs := &ast.ValueSpec{
			Names:  vs.Names,
			Type:   vs.Type,
			Values: vs.Values,
		}
		decl := &ast.GenDecl{
			Tok: token.CONST,
			Doc: gctx.renameComments(vs.Doc),
			Specs: []ast.Spec{
				newVs,
			},
		}
		outFile.Decls = append(outFile.Decls, decl)
	}

	sortedVars := sortSpecs(gctx.vars)
	for _, spec := range sortedVars {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		newVs := &ast.ValueSpec{
			Names:  vs.Names,
			Type:   vs.Type,
			Values: vs.Values,
		}
		decl := &ast.GenDecl{
			Tok: token.VAR,
			Doc: gctx.renameComments(vs.Doc),
			Specs: []ast.Spec{
				newVs,
			},
		}
		outFile.Decls = append(outFile.Decls, decl)
	}

	funcDecls := sortDecls(gctx.funcs)
	for _, decl := range funcDecls {
		fdecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		newFdecl := &ast.FuncDecl{
			Doc:  gctx.renameComments(fdecl.Doc),
			Recv: fdecl.Recv,
			Name: fdecl.Name,
			Type: fdecl.Type,
			Body: fdecl.Body,
		}
		outFile.Decls = append(outFile.Decls, newFdecl)
	}

	// newTokenPositioner().fixPositions(outFile)
	clearPositions(outFile)

	// ast.Print(outFset, outFile)

	// TODO:
	// If generating into different package than source, figure out dependencies,
	// copy private stuff, reference public stuff.

	buff := &bytes.Buffer{}

	err = printer.Fprint(buff, outFset, outFile)
	if err != nil {
		return errors.Wrap(err, "writing file failed")
	}

	outBytes, err := imports.Process(outFilename, buff.Bytes(), nil)
	if err != nil {
		return errors.Wrap(err, "Formatting file failed")
	}
	_, err = out.Write(outBytes)
	return errors.Wrap(err, "writing file failed")
}
