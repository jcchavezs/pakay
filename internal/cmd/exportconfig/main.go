package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/jcchavezs/pakay/internal/sources"
)

func main() {
	w := &bytes.Buffer{}

	fmt.Fprintln(w, "//go:generate go run ./internal/cmd/exportconfig")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "package pakay")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "import (")
	for _, s := range sources.GetAll() {
		e := reflect.ValueOf(s.ConfigFactory()).Elem()
		pkg := e.Type().PkgPath()
		fmt.Fprintf(w, "\t%q\n", pkg)
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "type (")
	for _, s := range sources.GetAll() {
		f := s.ConfigFactory()
		e := reflect.ValueOf(f).Elem()

		var prefix string
		switch f.Type() {
		case "1password":
			prefix = "OnePassword"
		default:
			prefix = f.Type()
		}

		fmt.Fprintf(w, "\t%s%sConfig = %s\n", strings.ToUpper(string(prefix[0])), prefix[1:], e.Type())
	}
	fmt.Fprintln(w, ")")

	if err := os.WriteFile("configexport.go", w.Bytes(), 0644); err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}
