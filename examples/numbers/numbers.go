package main

type Number int

const ZeroNumber Number = 0

var SomeNumber Number = 42

func AddNumber(a, b Number) Number {
	return a + b
}

func SubNumber(a, b Number) Number {
	return a - b
}
