package main

type StringIntSliceMap map[string][]int

func (m StringIntSliceMap) Flatten() []int {
	var ret []int
	for _, v := range m {
		ret = append(ret, v...)
	}
	return ret
}
