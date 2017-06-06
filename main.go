package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

const myName = "rei"

const (
	_ = iota
	exitcodeInvalidArgs
	exitcodeInvalidTypeMapping
	exitcodeDestFileFailed
	exitcodeSourceFileInvalid
	exitcodeGenFailed
)

func usage() {
	fmt.Fprintln(os.Stderr, `Usage: `+myName+` -in {source} [-out {dest}] "{types}"

Generates concrete code from generic code.

{source}  - (required) Source file with generic types
{dest}    - (optional) Destination file
{types}   - (required) Type mapping

Type mapping is in the following format:
  {generic1}={concrete1},[{generic2}={concrete2}]
where concrete can be one of the following:
  ConcreteType
`+"\t"+`concrete type in the same package
  pkg/pkg/pkg.ConcreteType
`+"\t"+`concrete type in a different package
  ("pkg/pkg/go-pkg")pkg.ConcreteType
`+"\t"+`concrete type in a different package, package name
`+"\t"+`doesn't match directory name

Flags:`)
	flag.PrintDefaults()
}

func fatal(code int, a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
	os.Exit(code)
}

func main() {
	var (
		in  = flag.String("in", "", "generic file")
		out = flag.String("out", "", "file to save output to instead of stdout")
	)
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		usage()
		os.Exit(exitcodeInvalidArgs)
	}

	if len(*in) == 0 {
		usage()
		os.Exit(exitcodeInvalidArgs)
	}

	typeMapping, err := parseMapping(args[0])
	if err != nil {
		fatal(exitcodeInvalidTypeMapping, err)
	}

	var file *os.File
	file, err = os.Open(*in)
	if err != nil {
		fatal(exitcodeSourceFileInvalid, err)
	}
	defer file.Close()

	// if targetPackageName is empty, gen will use the source package's name.
	targetPackageName := ""

	var outWriter io.Writer
	var outFilename string
	if len(*out) > 0 {
		targetDirName := path.Base(path.Dir(*out))
		if rel, err := filepath.Rel(path.Dir(*in), path.Dir(*out)); err != nil || rel != "." {
			// not the same directory, use directory name for generated package
			targetPackageName = targetDirName
		}
		err = os.MkdirAll(path.Dir(*out), 0755)
		if err != nil {
			fatal(exitcodeDestFileFailed, err)
		}

		outFile, err := os.Create(*out)
		if err != nil {
			fatal(exitcodeDestFileFailed, err)
		}
		defer outFile.Close()
		outWriter = outFile
		outFilename = *out
	} else {
		outWriter = os.Stdout
		outFilename = "stdout"
	}

	buffer := &bytes.Buffer{}

	err = gen(file, *in, targetPackageName, typeMapping, buffer, outFilename)
	if err != nil {
		fatal(exitcodeGenFailed, err)
	}

	_, err = io.Copy(outWriter, buffer)
	if err != nil {
		fatal(exitcodeGenFailed, err)
	}
}
