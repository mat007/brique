package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mat007/b"
)

func main() {
	defer building.CatchFailure()
	// $$$$ MAT: support some flags
	// -o build.exe
	// -f ./cmd/build
	os.Exit(run("./cmd/build", os.Args[1:]))
}

func run(dir string, args []string) int {
	b := building.Init("github.com/mat007/b")
	g := b.Go()
	buf := &bytes.Buffer{}
	g.WithOutput(buf).Run("list")
	path := strings.TrimSpace(buf.String())
	// $$$$ MAT: parse recursively?
	mainCode, pkgCode, isMain, err := parse(dir, path)
	if err != nil {
		building.Fatalln("parse failed:", err)
	}
	// If the user package is 'main' create the main file in the user folder
	// else put it in system temporary folder.
	todir := dir
	if !isMain {
		// $$$$ MAT: does not work in container build
		todir = os.TempDir()
	}
	uuid := strings.Replace(path, "/", "_", -1)
	// Relies on the fact that «The declaration order of variables declared in
	// multiple files is determined by the order in which the files are
	// presented to the compiler» and «To ensure reproducible initialization
	// behavior, build systems are encouraged to present multiple files
	// belonging to the same package in lexical file name order to a compiler»
	// See https://golang.org/ref/spec#Package_initialization
	pkgFile := filepath.Join(dir, "aaa_"+uuid+"_build.go")
	if err := ioutil.WriteFile(pkgFile, []byte(pkgCode), 0666); err != nil {
		building.Fatalln("write failed:", err)
	}
	defer os.Remove(pkgFile)
	mainFile := filepath.Join(todir, "aaa_"+uuid+"_main.go")
	if err := ioutil.WriteFile(mainFile, []byte(mainCode), 0666); err != nil {
		building.Fatalln("write failed:", err)
	}

	build := mainFile
	if isMain {
		defer os.Remove(mainFile)
		build = dir
	}
	g.Run("build", "-o", "build"+b.ExecExt(runtime.GOOS), build)
	return b.Command("build").WithSuccess().Run(args...)
}

type target struct {
	target string
	desc   string
	pkg    string
	name   string
}

func parse(dir, path string) (string, string, bool, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, parser.ParseComments)
	if err != nil {
		return "", "", false, err
	}
	var targets []target
	root := path + strings.TrimSuffix(strings.TrimPrefix(dir, "."), "/")
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				n, ok := decl.(*ast.FuncDecl)
				if !ok {
					continue
				}
				if n.Recv != nil {
					continue
				}
				name := n.Name.String()
				if isTarget(name, "Target") {
					err := checkFunc(fset, n, "B")
					if err != nil {
						return "", "", false, err
					}
					tgt, desc := makeTarget(name, n.Doc.Text())
					targets = append(targets, target{
						target: tgt,
						desc:   desc,
						pkg:    pkg.Name,
						name:   name,
					})
				}
			}
		}
	}
	if len(targets) == 0 {
		return "", "", false, fmt.Errorf("no targets found")
	}
	isMain := targets[0].pkg == "main"
	mainCode := `package main

import "github.com/mat007/b"
`
	if !isMain {
		mainCode += `import "` + root + `"
`
	}
	mainCode += `
func main() {
	defer building.CatchFailure()
	b := building.Builder()
`
	for _, t := range targets {
		name := t.name
		if !isMain {
			name = t.pkg + "." + t.name
		}
		mainCode += `	b.MakeTarget("` + t.target + `", "` + t.desc + `", ` + name + ")\n"
	}
	mainCode += `	b.Run()
}
`
	pkgCode := `package ` + targets[0].pkg + `

import "github.com/mat007/b"

var _ = building.Init("` + path + `")
`
	return mainCode, pkgCode, isMain, nil
}

// https://golang.org/pkg/cmd/go/internal/test/
func isTarget(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) {
		return true
	}
	c, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(c)
}

// https://golang.org/pkg/cmd/go/internal/test/
func checkFunc(fset *token.FileSet, fn *ast.FuncDecl, arg string) error {
	if !isFunc(fn, arg) {
		name := fn.Name.String()
		pos := fset.Position(fn.Pos())
		return fmt.Errorf("%s: wrong signature for %s, must be: func %s(%s *building.%s)",
			pos, name, name, strings.ToLower(arg), arg)
	}
	return nil
}

// https://golang.org/pkg/cmd/go/internal/test/
func isFunc(fn *ast.FuncDecl, arg string) bool {
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 ||
		fn.Type.Params.List == nil ||
		len(fn.Type.Params.List) != 1 ||
		len(fn.Type.Params.List[0].Names) > 1 {
		return false
	}
	ptr, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}
	// We can't easily check that the type is *testing.M
	// because we don't know how testing has been imported,
	// but at least check that it's *M or *something.M.
	// Same applies for B and T.
	if name, ok := ptr.X.(*ast.Ident); ok && name.Name == arg {
		return true
	}
	if sel, ok := ptr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == arg {
		return true
	}
	return false
}

func makeTarget(name, doc string) (string, string) {
	doc = strings.TrimPrefix(doc, name)
	crlf := strings.IndexAny(doc, "\r\n")
	if crlf != -1 {
		doc = doc[:crlf]
	}
	return makeTargetName(name), strings.TrimSpace(doc)
}

func makeTargetName(name string) string {
	name = strings.TrimPrefix(name, "Target")
	target := ""
	upper := true
	for _, c := range name {
		s := string(c)
		lower := strings.ToLower(s)
		if s != lower {
			if !upper {
				target += "-"
			}
			upper = true
		} else {
			upper = false
		}
		target += lower
	}
	return target
}
