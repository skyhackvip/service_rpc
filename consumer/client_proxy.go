package consumer

import (
	"errors"
	"net"
	"reflect"
	"strings"
)

func Call(funcName string, localFunc interface{}, params ...interface{}) (interface{}, error) {
	//get service
	service, err := split(funcName)
	if err != nil {
		return nil, err
	}

	//connect server
	conn, err := connHost(getHost(service.AppId))
	if err != nil {
		return nil, err
	}
	cli := NewClient(conn)

	//make func
	cli.Call(service.Method, localFunc)

	//reflect call
	return reflectCall(map[string]interface{}{service.Method: localFunc}, service.Method, params...)
}

type Service struct {
	AppId  string
	Method string
}

//demo: user.GetUser
func split(fun string) (Service, error) {
	arr := strings.Split(fun, ".")
	service := Service{}
	if len(arr) != 2 {
		return service, errors.New("fun name inlegal")
	}
	service.AppId = arr[0]
	service.Method = arr[1]
	return service, nil
}

//从注册中心拿host
func getHost(appId string) string {
	return "10.12.33.101:8811"
}

//conn host
func connHost(host string) (net.Conn, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func reflectCall(funcMap map[string]interface{}, name string, params ...interface{}) ([]reflect.Value, error) {
	f := reflect.ValueOf(funcMap[name]).Elem()
	if len(params) != f.Type().NumIn() {
		return nil, errors.New("params not adapted")
	}

	in := make([]reflect.Value, len(params))
	for idx, param := range params {
		in[idx] = reflect.ValueOf(param)
	}
	result := f.Call(in)
	return result, nil
}
