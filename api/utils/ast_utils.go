package utils

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func ParseGoFile(inputPath string) interface{} {
	codeMap := make(map[string]map[string]string)
	codeMap["package"] = make(map[string]string)
	codeMap["imports"] = make(map[string]string)
	codeMap["functions"] = make(map[string]string)
	codeMap["structs"] = make(map[string]string)

	fset := token.NewFileSet()
	src, err := os.ReadFile(inputPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	codeMap["package"][f.Name.Name] = "package " + f.Name.Name
	importStrings := []string{}

	for _, imp := range f.Imports {
		importStrings = append(importStrings, strings.TrimSpace(imp.Path.Value))
	}

	joinedImports := "import (\n    " + strings.Join(importStrings, "\n    ") + "\n)"
	codeMap["imports"]["all"] = joinedImports

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			// startLine := findCommentStart(fset, src, d.Pos())
			// endLine := fset.Position(d.End()).Line
			codeMap["functions"][d.Name.Name] = string(src[fset.Position(d.Pos()).Offset:fset.Position(d.End()).Offset])
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						codeMap["structs"][typeSpec.Name.Name] = string(src[fset.Position(typeSpec.Pos()).Offset:fset.Position(typeSpec.End()).Offset])
					}
				}
			}
		}
	}

	return codeMap
}

func findCommentStart(fset *token.FileSet, src []byte, pos token.Pos) int {
	line := fset.Position(pos).Line
	for line > 0 {
		lineStart := fset.Position(fset.File(pos).LineStart(line)).Offset
		lineEnd := fset.Position(fset.File(pos).LineStart(line+1)).Offset - 1
		if lineStart < 0 || lineEnd < 0 || lineEnd >= len(src) {
			break
		}
		lineStr := strings.TrimSpace(string(src[lineStart:lineEnd]))
		if lineStr == "" || strings.HasPrefix(lineStr, "//") {
			line--
		} else {
			break
		}
	}
	return line
}
