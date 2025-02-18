package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/google/go-jsonnet"
	"github.com/google/go-jsonnet/ast"
	importer "github.com/mashiike/go-jsonnet-alias-importer"
)

var nativeFunctions = []*jsonnet.NativeFunction{
	{
		Name:   "env",
		Params: []ast.Identifier{"name", "default"},
		Func: func(args []interface{}) (interface{}, error) {
			if len(args) != 2 {
				return nil, fmt.Errorf("env: invalid arguments length expected 2 got %d", len(args))
			}
			key, ok := args[0].(string)
			if !ok {
				return nil, fmt.Errorf("env: invalid 1st arguments, expected string got %T", args[0])
			}
			val := os.Getenv(key)
			if val == "" {
				return args[1], nil
			}
			return val, nil
		},
	},
}

//go:embed lib/*
var libFS embed.FS

func makeVM() *jsonnet.VM {
	vm := jsonnet.MakeVM()
	for _, f := range nativeFunctions {
		vm.NativeFunction(f)
	}
	im := importer.New()
	subFS, err := fs.Sub(libFS, "lib")
	if err != nil {
		log.Fatal(err)
	}
	im.Register("lib", subFS)
	vm.Importer(im)
	return vm
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("usage: main.go <jsonnet file>")
	}

	vm := makeVM()
	jsonStr, err := vm.EvaluateFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(jsonStr)
}
