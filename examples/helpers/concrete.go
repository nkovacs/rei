// Code generated by rei. DO NOT EDIT.

package main

func FooConcrete(a Concrete) Concrete {
	baz(a.ID)
	return barConcrete(a)
}
func barConcrete(a Concrete) Concrete {
	a.ID = 42
	return a
}
