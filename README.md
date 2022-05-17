# Peccary
GOLANG enthusiastic rest framework
It simply works :)


## Example
```go
package main

import (
	"fmt"
    "net/http"
	"github.com/xeuus/peccary"
)

type SomeService struct {
	counter uint64
}

type AuthService struct {
	token string
}

type HelloRequest struct {
	world string `json:"world"`
}

type Users struct {
}

func (u *Users) CreateUser(ctx *peccary.Context, rq *HelloRequest, svc *SomeService, au *AuthService) {
	svc.counter++
	ctx.EndWithData(http.StatusOK, svc.counter)
}

func (u *Users) Temp(ctx *peccary.Context, rq *Hello, svc *SomeService) {
	svc.counter++
	ctx.EndWithData(http.StatusOK, svc.counter)
}

func main() {
	api := peccary.NewPeccary()

	api.Singleton(&SomeService{})
	api.PerRequest(&AuthService{})

	api.Middleware(func(c *peccary.Context, svc *SomeService) {
		svc.counter++
		c.Next()
	})

	api.Controller(&Users{})
	api.Serve("/api", ":8080")
}

```
```bash
curl -s -X -v POST localhost:8080/api/users/createUser -d '{"world": ":)"}'
```
# Reference
## api := peccary.NewPeccary()
create a new Peccary instance
***
## api.Singleton(&Service{})
make the service globally available
***
## api.PerRequest(&Service{})
make the service available in request scope, in middlewares, controller. top to bottom.
***
## api.Proxy(&Service{})
make the service available per function call
***
## api.Middleware( func(*peccary.Context, *Service1, *Service2, ...) )
Middleware wraps the controller, can share functionality, have access to services, and simply can intercept user requests, \
Note: call context.Next() to call the next middleware or controller itself \
with having access to the context you have all the functionalities that a controller has
***
## api.Controller( &SomeController{} )
Controllers are simply http request handlers, you can call them via their name using curl:
Note: controller and method names interpreted in camel-case form
```
type SomeController struct {}
func (*SomeController) SayHello(ctx *peccary.Context, request string) {
    ctx.EndWithData(http.StatusOK, request)
}
api.Controller(&SomeController{})
curl localhost:8080/api/someController/sayHello -d '"Hello World"'
```