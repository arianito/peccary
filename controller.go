package peccary

import (
	"log"
	"reflect"
	"strings"
)

func (p *Peccary) Controller(ifc interface{}) {
	classType := reflect.TypeOf(ifc)
	classValue := reflect.ValueOf(ifc)

	classElem := classType
	if classElem.Kind() == reflect.Pointer {
		classElem = classElem.Elem()
	}

	for i := 0; i < classType.NumMethod(); i++ {
		m := classType.Method(i)

		var outputType reflect.Type
		methodType := m.Type
		hasOutput := false

		pth := []string{camelCase(classElem.Name()), camelCase(m.Name)}
		name := strings.ToLower(strings.Join(pth, "/"))

		for _, v := range p.handlers {
			if v.Name == name {
				log.Fatalf("Multiple controllers with same identifier found: %s", name)
			}

		}

		if methodType.NumIn() < 3 {
			log.Fatalf("Invalid number of arguments, Request parameter is required at %s(..)", name)
		}

		inputType := methodType.In(2)
		var args []reflect.Type

		if methodType.NumIn() > 3 {
			for j := 3; j < methodType.NumIn(); j++ {
				typ := methodType.In(j)
				args = append(args, typ)
			}
		}
		if methodType.NumOut() > 1 {
			hasOutput = true
			outputType = methodType.Out(0)
		}

		p.handlers = append(p.handlers, &Handler{
			Name:       name,
			Class:      classElem,
			Parent:     classValue,
			InputType:  inputType,
			HasOutput:  hasOutput,
			OutputType: outputType,
			Func:       m.Func,
			Args:       args,
		})
	}
}

func (p *Peccary) getController(name string) (*Handler, error) {
	for _, h := range p.handlers {
		if h.Name == name {
			return h, nil
		}
	}
	return nil, &EndpointNotFoundError{EndpointName: name}
}
