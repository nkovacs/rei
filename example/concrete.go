package main

import (
	models "github.com/nkovacs/rei/example/models"
)

type (
	ConcreteDAO       struct{}
	ConcreteDependant struct{ m models.Concrete }
)

var (
	zeroConcrete         models.Concrete
	ConcreteA, ConcreteB models.Concrete
)

func NewConcreteDAO() *ConcreteDAO {
	return &ConcreteDAO{}
}
func (dao *ConcreteDAO) Get(id int64) (*models.Concrete, error) {
	var m models.Concrete
	m.ID = id
	return &m, nil
}
func (dao *ConcreteDAO) Set(id int64) {
	type (
		Type struct{ IDD int64 }
	)
	var m Type
	m.IDD = id
}
func (dao *ConcreteDAO) GetType() {
}
func SomethingConcrete(d ConcreteDependant) {
}
func (dao *ConcreteDAO) Foo() {
	commonhelper()
}
func foobarConcreteA() models.Concrete {
	return ConcreteA
}
func foobarConcreteB() models.Concrete {
	return ConcreteB
}
func foobarConcreteBoth() (models.Concrete, models.Concrete) {
	return ConcreteA, ConcreteB
}
func foobarConcrete(t models.Concrete) {
}
func barbazConcrete(d ConcreteDependant) {
}
func barfooConcrete() {
	foobarConcrete(zeroConcrete)
}
