package main

import (
	"io/ioutil"

	test "github.com/nkovacs/rei/examples/pointer/go-test"
)

func ReadAllStringFromTestReader(r *test.TestReader) (string, error) {
	b, err := ioutil.ReadAll(r)
	return string(b), err
}
