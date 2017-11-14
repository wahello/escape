package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ankyra/escape/cmd"
	"github.com/spf13/cobra/doc"
)

const fmTemplate = `---
date: 2017-11-11 00:00:00
title: "%s"
slug: %s
type: "docs"
toc: true
---
`

const planHeader = `---
date: 2017-11-11 00:00:00
title: "The Escape Plan"
slug: escape-plan
type: "docs"
toc: true
---

The Escape Plan describes a package. 

Field | Type | Description
------|------|-------------
`

var typeMap = map[string][]string{
	"extends":          []string{"[string]", "Extensions"},
	"depends":          []string{"[string]", "Dependencies"},
	"consumes":         []string{"[string]", "Consumers"},
	"build_consumes":   []string{"[string]", "Consumers"},
	"deploy_consumes":  []string{"[string]", "Consumers"},
	"provides":         []string{"[string]", "Consumers"},
	"inputs":           []string{"[string]", "Variables"},
	"build_inputs":     []string{"[string]", "Variables"},
	"deploy_inputs":    []string{"[string]", "Variables"},
	"outputs":          []string{"[string]", "Variables"},
	"metadata":         []string{"{}"},
	"includes":         []string{"[]string"},
	"errands":          []string{"Errands"},
	"downloads":        []string{"Downloads"},
	"templates":        []string{"Templates"},
	"build_templates":  []string{"Templates"},
	"deploy_templates": []string{"Templates"},
}

var typeLinks = map[string]string{
	"Extensions":   "extensions",
	"Dependencies": "dependencies",
	"Consumers":    "providers-and-consumers",
	"Variables":    "input-and-output-variables",
	"Errands":      "errands",
	"Downloads":    "downloads",
	"Templates":    "templates",
}

type Page struct {
	Name       string
	Slug       string
	SrcFile    string
	StructName string
}

var Pages = map[string]Page{
	"escape plan": Page{"The Escape Plan", "escape-plan", "model/escape_plan/escape_plan.go", "EscapePlan"},
}

const PageHeader = `---
date: 2017-11-11 00:00:00
title: "%s"
slug: %s
type: "docs"
toc: true
---

%s

Field | Type | Description
------|------|-------------
%s
`

func GetYamlFieldFromTag(tag string) string {
	for _, s := range strings.Split(tag, " ") {
		s = strings.Trim(s, "`")
		if strings.HasPrefix(s, "yaml:\"") {
			s = s[6 : len(s)-1]
			parts := strings.Split(s, ",")
			return parts[0]
		}
	}
	return ""
}

func ParseType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return ParseType(t.X) + "." + t.Sel.String() // probably wrong
	case *ast.ArrayType:
		return "[" + ParseType(t.Elt) + "]"
	case *ast.StarExpr:
		return ParseType(t.X)
	case *ast.MapType:
		return "{" + ParseType(t.Key) + ":" + ParseType(t.Value) + "}"
	case *ast.InterfaceType:
		return "any"
	default:
		fmt.Printf("%T\n", t)
		panic("type not supported in documentation: ")
	}
	return ""
}

func StructTable(page Page, topLevelDoc string, s *ast.TypeSpec) string {
	structType := s.Type.(*ast.StructType)
	result := ""
	for _, field := range structType.Fields.List {
		tag := GetYamlFieldFromTag(field.Tag.Value)
		typ := ParseType(field.Type)
		result += "|" + tag + "|`" + typ + "`|"
		doc := strings.TrimSpace(field.Doc.Text())
		if doc != "" {
			for _, line := range strings.Split(doc, "\n") {
				if strings.HasPrefix(line, "#") {
					line = line[1:]
				}
				line = strings.TrimSpace(line)
				if line == "" {
					result += "\n|||"
				} else {
					result += line + " "
				}
			}
		}
		result += "\n"
	}
	return fmt.Sprintf(PageHeader, page.Name, page.Slug, topLevelDoc, result)
}

func GenerateStructDocs(f *ast.File, page Page) string {
	for _, decl := range f.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE {
			for _, spec := range gen.Specs {
				if s, ok := spec.(*ast.TypeSpec); ok {
					switch s.Type.(type) {
					case *ast.StructType:
						if s.Name.String() == page.StructName {
							return StructTable(page, gen.Doc.Text(), s)
						}
					}
				}
			}
		}
	}
	return ""
}

func GeneratePages() {
	for _, page := range Pages {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, page.SrcFile, nil, parser.ParseComments)
		if err != nil {
			panic(err)
		}
		str := GenerateStructDocs(f, page)
		filename := "docs/generated/" + page.Slug + ".md"
		fmt.Println("Writing ", filename)
		ioutil.WriteFile(filename, []byte(str), 0644)
	}
}

func main() {
	os.Mkdir("docs/generated/", 0755)
	filePrepender := func(filename string) string {
		name := filepath.Base(filename)
		base := strings.TrimSuffix(name, path.Ext(name))
		return fmt.Sprintf(fmTemplate, strings.Replace(base, "_", " ", -1), base)
	}

	linkHandler := func(name string) string {
		base := strings.TrimSuffix(name, path.Ext(name))
		return "../" + strings.ToLower(base) + "/"
	}
	err := doc.GenMarkdownTreeCustom(cmd.RootCmd, "./docs/cmd", filePrepender, linkHandler)
	if err != nil {
		log.Fatal(err)
	}
	GeneratePages()
}
