package provider

import (
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/data"
	"github.com/skyhackvip/service_rpc/transport"
	"io"
	"log"
	"net"
	"reflect"
)

type RPCServer struct {
	addr     string
	funcs    map[string]reflect.Value //keep service-func
	listener net.Listener
}

func NewRPCServer(addr string) *RPCServer {
	return &RPCServer{
		addr:  addr,
		funcs: make(map[string]reflect.Value), //route
	}
}

//register service
func (svr *RPCServer) Register(name string, function interface{}) {
	if _, ok := svr.funcs[name]; ok {
		log.Printf("%s is registered!\n", name)
		return
	}
	svr.funcs[name] = reflect.ValueOf(function)
	log.Printf("%s registered success!\n", name)
}

//service start
func (svr *RPCServer) Run() {
	//listen on port by tcp
	listener, err := net.Listen("tcp", svr.addr)
	if err != nil {
		panic(err)
	}
	log.Printf("listen on %s success!", svr.addr)
	svr.listener = listener

	//accept conn
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		//create new routine worker each connection
		go svr.handleConn(conn)
	}
}

//handle conn
func (svr *RPCServer) handleConn(conn net.Conn) {
	log.Println("---- handle conn ----")

	trans := transport.NewTransport(conn) //transport include conn and handle read/write

	rpcData := data.New(config.TRANS_TYPE) //data format:json/gob

	for {
		log.Println("---- loop start ----")

		//read from network
		reqData, err := trans.Read() //[]byte
		if err != nil {
			if err != io.EOF { //close
				log.Printf("read finish:%v\n", err)
				return
			}
			log.Printf("read panic:%v\n", err)
			return
		}
		log.Println("--- read data finish---")
		//decode data
		decodeData, err := rpcData.Decode(reqData) //rpcdata
		if err != nil {
			log.Println("decode request err:%v\n", err)
			return
		}
		log.Printf("decode data finish:%v\n", decodeData)

		//call local service
		result := svr.Call(decodeData)
		log.Println("call local service finish! result:", result)

		encodeRes, err := rpcData.Encode(result) //[]byte
		if err != nil {
			log.Printf("encode err:%v\n", err)
			return
		}
		log.Printf("encode result finish!")
		//send result
		err = trans.Send(encodeRes)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return
		}
		log.Printf("send result finish!")

		log.Println("---- loop end ----")
	}
}

//close tcp listener
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		log.Println("server close")
		svr.listener.Close()
	}
}

//call local service
func (svr *RPCServer) Call(req data.RPCData) data.RPCData {
	//get func
	f, ok := svr.funcs[req.Name]
	if !ok {
		log.Printf("%s is not exists\n", req.Name)
		return data.RPCData{}
	}

	//get arguments 如果有入参
	args := make([]reflect.Value, len(req.Args))
	for i := range req.Args {
		//reflect.TypeOf(req.Args[i])  //[]interface{}
		args[i] = reflect.ValueOf(req.Args[i]) //[1]
	}

	//start call
	result := f.Call(args) //reflect value类型，如果看value是func，可以调用Call方法执行

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
