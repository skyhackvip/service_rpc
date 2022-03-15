package provider

import (
	//	"github.com/skyhackvip/service_rpc/data"
	"reflect"
)

type Handler interface {
	Handle([]interface{}) ([]interface{}, error)
}

type RPCServerHandler struct {
	svr *RPCServer
	f   reflect.Value
}

//call local service
func (handler *RPCServerHandler) Handle(params []interface{}) ([]interface{}, error) {
	//get arguments 如果有入参
	args := make([]reflect.Value, len(params))
	for i := range params {
		//reflect.TypeOf(req.Args[i])  //[]interface{}
		args[i] = reflect.ValueOf(params[i]) //[1]
	}

	//start call
	result := handler.f.Call(args) //reflect value类型，如果看value是func，可以调用Call方法执行

	//result
	resArgs := make([]interface{}, len(result))
	for i := 0; i < len(result); i++ {
		resArgs[i] = result[i].Interface()
	}

	//error
	var err error
	if _, ok := result[len(result)-1].Interface().(error); ok {
		err = result[len(result)-1].Interface().(error)
	}

	return resArgs, err
}
