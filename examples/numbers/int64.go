package main

const ZeroInt64 int64 = 0

var SomeInt64 int64 = 42

func int64Adder(a, b int64) int64 {
	return a + b
}
func AddInt64(a, b int64) int64 {
	return int64Adder(a, b)
}
func SubInt64(a, b int64) int64 {
	return int64Adder(a, -b)
}
