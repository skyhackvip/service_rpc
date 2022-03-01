package consumer

import (
	"errors"
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/data"
	"github.com/skyhackvip/service_rpc/transport"
	"log"
	"net"
	"reflect"
)

type Client struct {
	conn net.Conn
}

func NewClient(conn net.Conn) *Client {
	return &Client{conn}
}

func (cli *Client) Call(name string, funcPtr interface{}) {
	log.Println("---- start call:", name)
	container := reflect.ValueOf(funcPtr).Elem() //反射获取函数元素

	f := func(req []reflect.Value) []reflect.Value {
		//出参个数
		numOut := container.Type().NumOut()

		//error
		errorHandler := func(err error) []reflect.Value {
			outArgs := make([]reflect.Value, numOut)
			for i := 0; i < len(outArgs)-1; i++ {
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
			outArgs[len(outArgs)-1] = reflect.ValueOf(&err).Elem()
			return outArgs
		}

		//conn
		trans := transport.NewTransport(cli.conn)

		//in args
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}
		reqRPC := data.RPCData{Name: name, Args: inArgs}
		log.Println("request RPC object", reqRPC)

		//encode request args
		rpcData := data.New(config.TRANS_TYPE)
		b, err := rpcData.Encode(reqRPC)
		if err != nil {
			log.Printf("encode err:%v\n", err)
			return errorHandler(err)
		}
		log.Println("encode success!")

		//send by network
		err = trans.Send(b)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return errorHandler(err)
		}
		log.Println("send success!")

		//get response
		resp, err := trans.Read() //[]byte
		if err != nil {
			return errorHandler(err)
		}
		log.Println("response success!")

		//decode response
		respDecode, err := rpcData.Decode(resp) //rpcdata
		if respDecode.Err != "" {
			log.Printf("decode err:%v\n", respDecode.Err)
			return errorHandler(errors.New(respDecode.Err))
		}
		if err != nil {
			log.Printf("decode err:%v\n", respDecode.Err)
			return errorHandler(err)
		}
		log.Println("decode success!", respDecode)

		//output result
		if len(respDecode.Args) == 0 {
			respDecode.Args = make([]interface{}, numOut)
		}
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut { //处理非error
				if respDecode.Args[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(respDecode.Args[i])
				}
			} else { //处理error
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}

	container.Set(reflect.MakeFunc(container.Type(), f)) //构造函数
}
