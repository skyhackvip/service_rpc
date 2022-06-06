package consumer

import (
	"errors"
	"strings"
)

type Service struct {
	Class  string
	Method string
}

//demo: User.GetUserById
func NewService(servicePath string) (*Service, error) {
	arr := strings.Split(servicePath, ".")
	service := &Service{}
	if len(arr) != 2 {
		return service, errors.New("service path inlegal")
	}
	service.Class = arr[0]
	service.Method = arr[1]
	return service, nil
}
