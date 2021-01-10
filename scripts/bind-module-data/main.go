package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

var modules = []string{
	"./sketch/core",
	"./sketch/stdlib/strings",
}

type Source struct {
	Filename string
	Code     string
}

type Variables struct {
	Package string
	Sources []*Source
}

var templ = template.Must(template.New("sketchCode").Parse(`
package {{ .Package }}

// Code generated by sketch/scripts/bind-module-data DO NOT EDIT
//
// Sources: {{ range .Sources }}
// {{ .Filename }}{{ end }}

const SketchCode = ` + "`" + `{{ range .Sources }}
; {{ .Filename }}
{{ .Code }}
{{ end }}` + "`\n"))

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	for _, module := range modules {
		if err := bindModuleData(module); err != nil {
			return err
		}
	}
	return nil
}

func bindModuleData(module string) error {

	sketchFiles, err := filepath.Glob(filepath.Join(module, "*.skt"))
	if err != nil {
		return err
	}

	var sources []*Source
	for _, sketchFile := range sketchFiles {
		// Format Sketch files before reading
		sketchFmtCmd := exec.Command("sketch", "format", "-w", sketchFile)
		if err := sketchFmtCmd.Run(); err != nil {
			return err
		}

		data, err := ioutil.ReadFile(sketchFile)
		if err != nil {
			return err
		}

		source := &Source{
			Filename: filepath.Base(sketchFile),
			Code:     strings.TrimSpace(string(data)),
		}

		sources = append(sources, source)
	}

	outputFile := filepath.Join(module, "sketch_code.go")
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer f.Close()

	variables := &Variables{
		Package: filepath.Base(module),
		Sources: sources,
	}

	if err := templ.Execute(f, variables); err != nil {
		return err
	}

	// Format new Go file
	goFmtCmd := exec.Command("gofmt", "-w", outputFile)
	if err := goFmtCmd.Run(); err != nil {
		return err
	}

	return nil
}
