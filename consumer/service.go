package consumer

import (
	"errors"
	"math/rand"
	"strings"
	"time"
)

type Service struct {
	AppId  string
	Class  string
	Method string
	Addrs  []string
}

//demo: UserService.User.GetUserById
func NewService(servicePath string) (*Service, error) {
	arr := strings.Split(servicePath, ".")
	service := &Service{}
	if len(arr) != 3 {
		return service, errors.New("service path inlegal")
	}
	service.AppId = arr[0]
	service.Class = arr[1]
	service.Method = arr[2]
	return service, nil
}

//selector: random rb weight
func (service *Service) SelectAddr() string {
	rand.Seed(time.Now().Unix())
	return service.Addrs[rand.Intn(len(service.Addrs))]
}
