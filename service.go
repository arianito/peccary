package peccary

import "reflect"

const (
	singleton = 0
	request   = 1
	proxy     = 2
)

type service struct {
	name     string
	mode     uint8
	class    reflect.Type
	instance reflect.Value
}

func (p *Peccary) registerService(ifc interface{}, mode uint8) {
	classType := reflect.TypeOf(ifc)
	classValue := reflect.ValueOf(ifc)

	classElem := classType
	if classElem.Kind() == reflect.Pointer {
		classElem = classElem.Elem()
	}

	name := classElem.Name()

	p.services = append(p.services, &service{
		name:     name,
		class:    classElem,
		instance: classValue,
		mode:     mode,
	})
}
func (p *Peccary) Singleton(ifc interface{}) {
	p.registerService(ifc, singleton)
}
func (p *Peccary) PerRequest(ifc interface{}) {
	p.registerService(ifc, request)
}
func (p *Peccary) Proxy(ifc interface{}) {
	p.registerService(ifc, proxy)
}

func (p *Peccary) getService(c *Context, arg reflect.Type) (reflect.Value, error) {
	isPtr := arg.Kind() == reflect.Pointer
	argElem := arg
	if isPtr {
		argElem = argElem.Elem()
	}
	for _, svc := range p.services {
		if svc.name == argElem.Name() {
			var instance reflect.Value
			if svc.mode == request {
				cachedInstance, ok := c.services[svc.name]
				if ok {
					instance = cachedInstance
				} else {
					instance = reflect.New(svc.class)
					c.services[svc.name] = instance
				}
			} else if svc.mode == proxy {
				instance = reflect.New(svc.class)
			} else if svc.mode == singleton {
				instance = svc.instance
			}
			if !isPtr {
				instance = instance.Elem()
			}
			return instance, nil
		}
	}
	return reflect.Value{}, &ServiceNotFoundError{ServiceName: argElem.Name()}
}
