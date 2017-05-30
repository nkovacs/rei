package main

type Type interface{}
type TypeSlice []Type

func (s TypeSlice) Where(fn func(Type) bool) (result TypeSlice) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}
