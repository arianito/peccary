package peccary

import (
	"fmt"
)

type ServiceNotFoundError struct {
	ServiceName string
}

func (e *ServiceNotFoundError) Error() string {
	return fmt.Sprintf("%s not registered", e.ServiceName)
}

type EndpointNotFoundError struct {
	EndpointName string
}

func (e *EndpointNotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.EndpointName)
}
