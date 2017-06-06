package main

import (
	"io/ioutil"
	"os"
)

func ReadAllStringFromFile(r *os.File) (string, error) {
	b, err := ioutil.ReadAll(r)
	return string(b), err
}
