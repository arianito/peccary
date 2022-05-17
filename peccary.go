package peccary

import (
	"reflect"
)

type Handler struct {
	Name       string
	Class      reflect.Type
	Parent     reflect.Value
	InputType  reflect.Type
	HasOutput  bool
	OutputType reflect.Type
	Func       reflect.Value
	Args       []reflect.Type
}

type Peccary struct {
	services    []*service
	handlers    []*Handler
	middlewares []*middleware
}

func NewPeccary() *Peccary {
	return &Peccary{}
}
