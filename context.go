package peccary

import (
	"net/http"
	"reflect"
)

type Context struct {
	Writer      http.ResponseWriter
	Request     *http.Request
	Storage     map[string]interface{}
	canContinue bool
	chainIndex  uint8
	chainArgs   [][]reflect.Value
	chainFunc   []reflect.Value
	services    map[string]reflect.Value
	statusCode  int
	err         error
	output      interface{}
}

func (c *Context) SetStatusCode(statusCode int) {
	c.statusCode = statusCode
}
func (c *Context) SetData(data interface{}) {
	c.output = data
}
func (c *Context) SetError(err error) {
	c.err = err
}

func (c *Context) End() {
	c.canContinue = false
}
func (c *Context) EndWithData(statusCode int, data interface{}) {
	c.SetStatusCode(statusCode)
	c.SetData(data)
	c.End()
}

func (c *Context) EndWithStatus(statusCode int) {
	c.SetStatusCode(statusCode)
	c.End()
}

func (c *Context) EndWithError(statusCode int, err error) {
	c.SetStatusCode(statusCode)
	c.SetError(err)
	c.End()
}

func (c *Context) Next() {
	if !c.canContinue {
		return
	}
	index := c.chainIndex
	c.chainIndex++
	c.chainFunc[index].Call(c.chainArgs[index])
}
