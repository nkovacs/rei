package main

type Number int

const ZeroNumber Number = 0

var SomeNumber Number = 42

func numberAdder(a, b Number) Number {
	return a + b
}

func AddNumber(a, b Number) Number {
	return numberAdder(a, b)
}

func SubNumber(a, b Number) Number {
	return numberAdder(a, -b)
}
