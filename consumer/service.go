package consumer

import (
	"errors"
	"strings"
)

type Service struct {
	AppId  string
	Method string
	Addrs  []string
}

//demo: test.HelloWorld user.GetUser
func NewService(methodName string) (*Service, error) {
	arr := strings.Split(methodName, ".")
	service := &Service{}
	if len(arr) != 2 {
		return service, errors.New("method name inlegal")
	}
	service.AppId = arr[0]
	service.Method = arr[1]
	return service, nil
}

func (service *Service) SelectAddr() string {
	return "10.12.33.101:8811"
}
