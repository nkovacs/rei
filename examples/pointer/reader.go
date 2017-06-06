package main

import (
	"io"
	"io/ioutil"
)

type Reader interface {
	io.Reader
}

func ReadAllStringFromReader(r Reader) (string, error) {
	b, err := ioutil.ReadAll(r)
	return string(b), err
}
