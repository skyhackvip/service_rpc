package provider

import (
	"github.com/skyhackvip/service_rpc/data"
	"reflect"
)

type Handler interface {
	Handle(req data.RPCData) data.RPCData
}

type RPCServerHandler struct {
	svr *RPCServer
	f   reflect.Value
}

//call local service
func (handler *RPCServerHandler) Handle(req data.RPCData) data.RPCData {
	//get func
	/*f, ok := svr.funcs[req.Name]
	if !ok {
		log.Printf("%s is not exists\n", req.Name)
		return data.RPCData{}
	}*/

	//get arguments 如果有入参
	args := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		//reflect.TypeOf(req.Args[i])  //[]interface{}
		args[i] = reflect.ValueOf(req.Args[i]) //[1]
	}

	//start call
	result := handler.f.Call(args) //reflect value类型，如果看value是func，可以调用Call方法执行

	//result
	resArgs := make([]interface{}, len(result))
	for i := 0; i < len(result); i++ {
		resArgs[i] = result[i].Interface()
	}

	//error
	var err string
	if _, ok := result[len(result)-1].Interface().(error); ok {
		err = result[len(result)-1].Interface().(error).Error()
	}

	return data.RPCData{Name: req.Name, Args: resArgs, Err: err}
}
