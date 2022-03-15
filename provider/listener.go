package provider

import (
	"fmt"
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/global"
	"github.com/skyhackvip/service_rpc/protocol"
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
	//	closed      chan struct{}
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
	nl, err := net.Listen(config.NET_TRANS_PROTOCOL, addr)
	if err != nil {
		panic(err)
	}
	l.nl = nl
	log.Printf("listen on %s success!", addr)

	//accept conn
	for {
		conn, err := l.nl.Accept()
		if err != nil {
			//log.Printf("accept err: %v\n", err)
			continue
		}

		//create new routine worker each connection
		go l.handleConn(conn)
	}
}

//handle each connection
//TODO:对异常 err 处理
func (l *RPCListener) handleConn(conn net.Conn) {
	defer catchPanic()

	for {
		//read from network
		msg, err := l.receiveData(conn)
		if err != nil || msg == nil {
			return
		}

		//decode
		coder := global.Codecs[msg.Header.SerializeType()] //get from cache
		if coder == nil {
			return
		}
		inArgs := make([]interface{}, 0)
		err = coder.Decode(msg.Payload, &inArgs) //rpcdata
		if err != nil {
			log.Println("decode request err:%v\n", err)
			return
		}
		log.Printf("decode data finish:%v\n", inArgs)

		//call local service
		handler, ok := l.Handlers[msg.ServiceClass]
		if !ok {
			log.Println("can not found handler")
			return
		}
		result, err := handler.Handle(msg.ServiceMethod, inArgs)
		log.Println("call local service finish! result:", result)

		//encode
		encodeRes, err := coder.Encode(result) //[]byte result + err
		if err != nil {
			log.Printf("encode err:%v\n", err)
			return
		}

		//send result
		err = l.sendData(conn, encodeRes)
		if err != nil {
			log.Printf("send err:%v\n", err)
			return
		}
		log.Printf("send result finish!")
	}
}

func (l *RPCListener) receiveData(conn net.Conn) (*protocol.RPCMsg, error) {
	msg, err := protocol.Read(conn)
	if err != nil {
		if err != io.EOF { //close
			return nil, err
		}
		log.Printf("read finish:%v\n", err)
	}
	return msg, nil
}

func (l *RPCListener) sendData(conn net.Conn, payload []byte) error {
	resMsg := protocol.NewRPCMsg()
	resMsg.SetVersion(config.Protocol_MsgVersion)
	resMsg.SetMsgType(protocol.Response)
	resMsg.SetCompressType(protocol.None)
	resMsg.SetSerializeType(protocol.Gob)
	resMsg.Payload = payload
	_, err := resMsg.Send(conn)
	return err
}

func (l *RPCListener) Close() {
	if l.nl != nil {
		l.nl.Close()
	}
}

func catchPanic() {
	err := recover()
	if err != nil {
		log.Println("catch panic err:", err)
	}
}
