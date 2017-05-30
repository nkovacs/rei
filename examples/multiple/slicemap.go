package main

type KeyType interface{}
type ValueType interface{}
type KeyTypeValueTypeSliceMap map[KeyType][]ValueType

func (m KeyTypeValueTypeSliceMap) Flatten() []ValueType {
	var ret []ValueType
	for _, v := range m {
		ret = append(ret, v...)
	}
	return ret
}
