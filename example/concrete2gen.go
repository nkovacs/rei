package main

type (
	Concrete2DAO       struct{}
	Concrete2Dependant struct{ m Concrete2 }
)

var (
	zeroConcrete2          Concrete2
	Concrete2A, Concrete2B Concrete2
)

func NewConcrete2DAO() *Concrete2DAO {
	return &Concrete2DAO{}
}
func (dao *Concrete2DAO) Get(id int64) (*Concrete2, error) {
	var m Concrete2
	m.ID = id
	return &m, nil
}
func (dao *Concrete2DAO) Set(id int64) {
	type (
		Type struct{ IDD int64 }
	)
	var m Type
	m.IDD = id
}
func (dao *Concrete2DAO) GetType() {
}
func SomethingConcrete2(d Concrete2Dependant) {
}
func (dao *Concrete2DAO) Foo() {
	commonhelper()
}
func foobarConcrete2A() Concrete2 {
	return Concrete2A
}
func foobarConcrete2B() Concrete2 {
	return Concrete2B
}
func foobarConcrete2Both() (Concrete2, Concrete2) {
	return Concrete2A, Concrete2B
}
func foobarConcrete2(t Concrete2) {
}
func barbazConcrete2(d Concrete2Dependant) {
}
func barfooConcrete2() {
	foobarConcrete2(zeroConcrete2)
}
