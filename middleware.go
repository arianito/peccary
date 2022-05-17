package peccary

import (
	"log"
	"reflect"
)

type middleware struct {
	Func reflect.Value
	Args []reflect.Type
}

func (p *Peccary) Middleware(mdl interface{}) {
	m := reflect.ValueOf(mdl)
	methodType := m.Type()
	if methodType.NumIn() == 0 {
		log.Fatalf("Invalid number of arguments, *Context is required.")
	}
	var args []reflect.Type
	if methodType.NumIn() > 0 {
		for j := 1; j < methodType.NumIn(); j++ {
			typ := methodType.In(j)
			args = append(args, typ)
		}
	}
	p.middlewares = append(p.middlewares, &middleware{
		Args: args,
		Func: m,
	})
}
