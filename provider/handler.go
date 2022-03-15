package provider

import (
	"reflect"
)

type Handler interface {
	Handle(string, []interface{}) ([]interface{}, error)
}

type RPCServerHandler struct {
	svr   *RPCServer
	class reflect.Value
}

//call local service
func (handler *RPCServerHandler) Handle(method string, params []interface{}) ([]interface{}, error) {
	//get arguments if params is not empty
	args := make([]reflect.Value, len(params))
	for i := range params {
		args[i] = reflect.ValueOf(params[i]) //[1]
	}

	//get method
	reflectMethod := handler.class.MethodByName(method)

	result := reflectMethod.Call(args)

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
