package peccary

import (
	"encoding/json"
	"net/http"
	"reflect"
)

type Response struct {
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
	Status int         `json:"status,omitempty"`
}

func writeJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

type worker func(http.ResponseWriter, *http.Request)

func getControllerNameFromUrl(url string, prefix string) (string, error) {
	if url[:len(prefix)] != prefix {
		return "", &EndpointNotFoundError{EndpointName: url}
	}
	return url[len(prefix)+1:], nil
}

func (h worker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}
func (p *Peccary) Serve(prefix string, addr string) {
	var wrk worker = func(w http.ResponseWriter, r *http.Request) {
		var url string = r.URL.Path

		controllerName, err := getControllerNameFromUrl(url, prefix)

		if err != nil {
			writeJson(w, http.StatusBadRequest, &Response{
				Error: err.Error(),
			})
			return
		}

		var ctx *Context = &Context{
			Writer:      w,
			Request:     r,
			statusCode:  http.StatusOK,
			services:    make(map[string]reflect.Value),
			canContinue: true,
			chainIndex:  0,
		}

		handler, err := p.getController(controllerName)
		if err != nil {
			writeJson(w, http.StatusBadRequest, &Response{
				Error: err.Error(),
			})
			return
		}
		var chainArgs [][]reflect.Value
		var chainFunc []reflect.Value
		for _, m := range p.middlewares {
			arguments := []reflect.Value{reflect.ValueOf(ctx)}
			for _, arg := range m.Args {
				instance, err := p.getService(ctx, arg)
				if err != nil {
					writeJson(w, http.StatusInternalServerError, &Response{
						Error: err.Error(),
					})
					return
				}
				arguments = append(arguments, instance)
			}
			chainArgs = append(chainArgs, arguments)
			chainFunc = append(chainFunc, m.Func)
		}
		chainArgs = append(chainArgs, []reflect.Value{reflect.ValueOf(ctx)})
		chainFunc = append(chainFunc, reflect.ValueOf(func(ctx *Context) {
			typ := handler.InputType
			isPtr := typ.Kind() == reflect.Pointer
			if isPtr {
				typ = typ.Elem()
			}
			input := reflect.New(typ)
			decoder := json.NewDecoder(r.Body)
			if err := decoder.Decode(input.Interface()); err != nil {
				ctx.EndWithError(http.StatusBadRequest, err)
				return
			}
			if !isPtr {
				input = input.Elem()
			}
			arguments := []reflect.Value{handler.Parent, reflect.ValueOf(ctx), input}
			for _, arg := range handler.Args {
				instance, err := p.getService(ctx, arg)
				if err != nil {
					ctx.EndWithError(http.StatusInternalServerError, err)
					return
				}
				arguments = append(arguments, instance)
			}
			response := handler.Func.Call(arguments)
			rs := len(response)
			if rs > 0 {
				if err, ok := response[0].Interface().(error); ok && err != nil {
					ctx.EndWithError(http.StatusInternalServerError, err)
					return
				}
			}
			if rs > 1 {
				ctx.SetData(response[0].Interface())
				if err, ok := response[1].Interface().(error); ok && err != nil {
					ctx.EndWithError(http.StatusInternalServerError, err)
					return
				}
			}
		}))
		ctx.chainArgs = chainArgs
		ctx.chainFunc = chainFunc
		ctx.Next()
		response := &Response{
			Data: ctx.output,
		}
		if ctx.err != nil {
			response.Error = ctx.err.Error()
		}
		writeJson(w, ctx.statusCode, response)
	}
	http.ListenAndServe(addr, wrk)
}
