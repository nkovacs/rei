package main

type IntSlice []int

func (s IntSlice) Where(fn func(int) bool) (result IntSlice) {
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return
}
