package provider

import (
	"fmt"
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/data"
	"github.com/skyhackvip/service_rpc/transport"
	"io"
	"log"
	"net"
)

type Listener interface {
	Run()
	SetHandler(string, Handler)
	Close()
}

//base on tcp
type RPCListener struct {
	ServiceIp   string
	ServicePort int
	Handlers    map[string]Handler
	nl          net.Listener
}

func NewRPCListener(serviceIp string, servicePort int) *RPCListener {
	return &RPCListener{ServiceIp: serviceIp,
		ServicePort: servicePort,
		Handlers:    make(map[string]Handler)}
}

func (l *RPCListener) SetHandler(name string, handler Handler) {
	if _, ok := l.Handlers[name]; ok {
		log.Printf("%s is registered!\n", name)
		return
	}
	l.Handlers[name] = handler
}

//start listening and waiting for connection
func (l *RPCListener) Run() {
	//listen on port by tcp
	addr := fmt.Sprintf("%s:%d", l.ServiceIp, l.ServicePort)
	nl, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	l.nl = nl
	log.Printf("listen on %s success!", addr)

	//accept conn
	for {
		conn, err := l.nl.Accept()
		if err != nil {
			log.Printf("accept err: %v\n", err)
			continue
		}

		//create new routine worker each connection
		go l.handleConn(conn)
	}

}

//handle each connection
func (l *RPCListener) handleConn(conn net.Conn) {
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
		//get handler
		handler, ok := l.Handlers[decodeData.Name]
		if !ok {
			log.Println("can not found handler")
			return
		}

		result := handler.Handle(decodeData)
		log.Println("call local service finish! result:", result)

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

func (l *RPCListener) Close() {
	l.nl.Close()
}
