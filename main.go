package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cfgPath = kingpin.Flag("config", `path to JSON configuration, defaults to "config.json"`).Short('c').String()
	outFile = kingpin.Flag("out", `path to output file, defaults to "config.go"`).Short('o').String()
	pack    = kingpin.Flag("package", `name of go package containing config code`).Short('p').Default("main").String()

	// Replace the variables below for testing
	fatalf = kingpin.Fatalf
)

func main() {
	kingpin.Parse()
	run()
}

// Make it easy for test to setup configuration values.
func run() {
	wd, err := os.Getwd()
	if err != nil {
		fatalf("failed to retrieve current directory: %s", err)
	}
	if *cfgPath == "" {
		*cfgPath = path.Join(wd, "config.json")
	}
	if *cfgPath, err = filepath.Abs(*cfgPath); err != nil {
		fatalf("failed to compute absolute file path for config: %s", err)
	}
	if *outFile == "" {
		*outFile = path.Join(wd, "config.go")
	}
	if _, err := os.Stat(*cfgPath); os.IsNotExist(err) {
		fatalf("no configuration file at %s", *cfgPath)
	}
	outDir := filepath.Dir(*outFile)
	s, err := os.Stat(outDir)
	if os.IsNotExist(err) {
		kingpin.FatalIfError(os.MkdirAll(outDir, 0777), "output directory")
	} else if err != nil || !s.IsDir() {
		fatalf("not a valid output directory: %s", outDir)
	}
	input, err := os.Open(*cfgPath)
	if err != nil {
		fatalf("failed to open JSON file: %s", err)
	}
	defer input.Close()
	decoder := json.NewDecoder(input)
	decoder.UseNumber()
	var data interface{}
	err = decoder.Decode(&data)
	if err != nil {
		fatalf("failed to unmarshal JSON: %s", err)
	}
	var tree Tree
	tree.Populate(data)
	if tree.Type != Struct {
		fatalf("invalid configuration file content, JSON must define an object")
	}
	tree.Normalize()
	tree.Name = "cfg"
	configs := make([]*Tree, len(tree.Children))
	for i, child := range tree.Children {
		configs[i] = &Tree{
			Name:     child.Name + "Cfg",
			Children: child.Children,
			Type:     child.Type,
			List:     child.List,
		}
	}
	vars := Variables{
		CmdLine: fmt.Sprintf("$ %s %s", os.Args[0], strings.Join(os.Args[1:], " ")),
		Pack:    *pack,
		Tree:    &tree,
		Configs: configs,
	}
	tmpl, err := template.New("gonfig").Funcs(template.FuncMap{"NewVar": NewVar}).Parse(configTemplate)
	if err != nil {
		fatalf("failed to compile template: %s", err)
	}
	outF, err := os.Create(*outFile)
	if err != nil {
		fatalf("failed to create file: %s", err)
	}
	defer outF.Close()
	err = tmpl.Execute(outF, &vars)

	// Print something
	fmt.Println(*outFile)
}

// Variables defines the data structure fed to the template to generate the final code.
type Variables struct {
	CmdLine string  // Command line used to invoke gonfig
	Pack    string  // Name of target package
	Tree    *Tree   // Top level data structure
	Configs []*Tree // Configuration entries
}

// Internal count used to generate unique variable names
var varCount int

// Generate unique Go variable name
func NewVar() string {
	varCount += 1
	return fmt.Sprintf("v%d", varCount)
}

const configTemplate string = `//************************************************************************//
//                     Configuration
//
// Generated with:
// {{.CmdLine}}
//
// The content of this file is auto-generated, DO NOT MODIFY
//************************************************************************//

package {{.Pack}}

import (
	"encoding/json"
        "os"
)

var ({{range .Tree.Children}}
	{{.Name}} {{.TypeName}}{{end}}
)
{{range .Configs}}{{if eq .Type.String "struct"}}
{{.FormatRaw}}{{end}}{{end}}
{{.Tree.FormatRaw}}
// Load reads the JSON at the given path and initializes the package variables with the
// corresponding values.
func Load(path string) error {
        input, err := os.Open(path)
        if err != nil {
                return err
        }
	decoder := json.NewDecoder(input)
        var c Cfg
        err = decoder.Decode(&c)
        if err != nil {
                return err
	}{{range .Tree.Children}}{{$v := NewVar}}{{if eq .Type.String "struct"}}
	{{$v}} := {{.Name}}Cfg(c.{{.Name}}){{end}}
	{{.Name}} = {{if eq .Type.String "struct"}}&{{$v}}{{else}}c.{{.Name}}{{end}}{{end}}

        return nil
}
`

// Wrapper around tree.Format that exits the process in case of error.
func (t *Tree) FormatRaw() string {
	b, err := t.Format()
	if err != nil {
		fatalf("failed to format data structure: %s", err)
	}
	return string(b)
}

// Produce type name for tree node
func (t *Tree) TypeName() string {
	if t.Type == Struct {
		return "*" + t.Name.String() + "Cfg"
	}
	return t.Type.String()
}
