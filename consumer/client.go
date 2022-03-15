package consumer

import (
	"context"
	"errors"
	"fmt"
	"github.com/skyhackvip/service_rpc/codec"
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/protocol"
	"log"
	"net"
	"reflect"
	"time"
)

type Client interface {
	Connect(string) error
	Invoke(context.Context, string, interface{}, ...interface{}) (interface{}, error)
	Close()
}

type Option struct {
	Retries           int
	ConnectionTimeout time.Duration
	SerializeType     protocol.SerializeType
	CompressType      protocol.CompressType
}

var DefaultOption = Option{
	Retries:           3,
	ConnectionTimeout: 60 * time.Second,
	SerializeType:     protocol.Gob,
	CompressType:      protocol.None,
}

type RPCClient struct {
	conn   net.Conn
	option Option
}

func NewClient(option Option) Client {
	return &RPCClient{option: option}
}

func (cli *RPCClient) Connect(addr string) error {
	conn, err := net.DialTimeout(config.NET_TRANS_PROTOCOL, addr, cli.option.ConnectionTimeout)
	if err != nil {
		return err
	}
	cli.conn = conn
	return nil
}

func (cli *RPCClient) Invoke(ctx context.Context, methodName string, stub interface{}, params ...interface{}) (interface{}, error) {

	//make func : this step can be prepared before invoke and store into cache
	cli.makeCall(methodName, stub)

	//reflect call
	return cli.wrapCall(ctx, stub, params...)
}

func (cli *RPCClient) Close() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

//make call func
func (cli *RPCClient) makeCall(methodName string, methodPtr interface{}) {
	log.Println("---- start call:", methodName)
	container := reflect.ValueOf(methodPtr).Elem() //反射获取函数元素
	coder := codec.New(config.CODEC_GOB)

	handler := func(req []reflect.Value) []reflect.Value {
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

		//in args
		inArgs := make([]interface{}, 0, len(req))
		for _, arg := range req {
			inArgs = append(inArgs, arg.Interface())
		}

		payload, err := coder.Encode(inArgs) //[]byte
		if err != nil {
			log.Printf("encode err:%v\n", err)
			return errorHandler(err)
		}
		log.Println("encode success!", payload)

		//send by network
		msg := protocol.NewRPCMsg()
		msg.SetVersion(config.Protocol_MsgVersion)
		msg.SetMsgType(protocol.Request)
		msg.SetCompressType(cli.option.CompressType)
		msg.SetSerializeType(cli.option.SerializeType)
		msg.ServiceMethod = methodName
		msg.Payload = payload
		_, err = msg.Send(cli.conn)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return errorHandler(err)
		}
		log.Println("send success!")

		//read from network
		respMsg, err := protocol.Read(cli.conn)
		if err != nil {
			return errorHandler(err)
		}
		log.Println("response success!")

		//decode response
		respDecode := make([]interface{}, 0)
		err = coder.Decode(respMsg.Payload, &respDecode)
		if err != nil {
			log.Printf("decode err:%v\n", err)
			return errorHandler(err)
		}
		log.Println("decode success!", respDecode)

		//output result
		if len(respDecode) == 0 {
			respDecode = make([]interface{}, numOut)
		}
		outArgs := make([]reflect.Value, numOut)
		for i := 0; i < numOut; i++ {
			if i != numOut { //处理非error
				if respDecode[i] == nil {
					outArgs[i] = reflect.Zero(container.Type().Out(i))
				} else {
					outArgs[i] = reflect.ValueOf(respDecode[i])
				}
			} else { //处理error
				outArgs[i] = reflect.Zero(container.Type().Out(i))
			}
		}
		return outArgs
	}

	container.Set(reflect.MakeFunc(container.Type(), handler)) //构造函数
}

func (cli *RPCClient) wrapCall(ctx context.Context, stub interface{}, params ...interface{}) (interface{}, error) {
	f := reflect.ValueOf(stub).Elem()
	if len(params) != f.Type().NumIn() {
		return nil, errors.New(fmt.Sprintf("params not adapted: %d-%d", len(params), f.Type().NumIn()))
	}

	in := make([]reflect.Value, len(params))
	for idx, param := range params {
		in[idx] = reflect.ValueOf(param)
	}
	result := f.Call(in)
	return result, nil
}
