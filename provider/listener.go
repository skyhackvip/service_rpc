package provider

import (
	"fmt"
	"github.com/skyhackvip/service_rpc/config"
	"github.com/skyhackvip/service_rpc/global"
	"github.com/skyhackvip/service_rpc/protocol"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"
)

type Listener interface {
	Run()
	SetHandler(string, Handler)
	SetPlugins(PluginContainer)
	Close()
	GetAddrs() []string
	Shutdown()
}

//base on tcp
type RPCListener struct {
	ServiceIp   string
	ServicePort int
	option      Option
	Plugins     PluginContainer
	Handlers    map[string]Handler
	nl          net.Listener
	doneChan    chan struct{} //控制结束
	handlingNum int32         //处理中任务数
	shutdown    int32         //关闭处理中
}

func NewRPCListener(option Option) *RPCListener {
	return &RPCListener{ServiceIp: option.Ip,
		ServicePort: option.Port,
		option:      option,
		Handlers:    make(map[string]Handler),
		doneChan:    make(chan struct{}),
	}
}

func (l *RPCListener) SetPlugins(plugins PluginContainer) {
	l.Plugins = plugins
}

func (l *RPCListener) SetHandler(name string, handler Handler) {
	if _, ok := l.Handlers[name]; ok {
		log.Printf("%s is registered!\n", name)
		return
	}
	l.Handlers[name] = handler
}

//start listening and waiting for connection
func (l *RPCListener) Run() { //return
	//listen on port by tcp TODO
	addr := fmt.Sprintf("%s:%d", l.ServiceIp, l.ServicePort)
	log.Println(l.option.NetProtocol, addr)
	nl, err := net.Listen(l.option.NetProtocol, addr)
	if err != nil {
		panic(err)
	}
	//lock ?
	l.nl = nl
	log.Printf("listen on %s success!", addr)

	//accept conn
	for {
		conn, err := l.nl.Accept()
		if err != nil {
			select { //done
			case <-l.getDoneChan():
				log.Println("server closed done")
				return
			default:
			}

			if e, ok := err.(net.Error); ok && e.Temporary() { //网络发生临时错误
				log.Printf("server accept error: %v", err)
				time.Sleep(5 * time.Millisecond)
				continue
			}

			log.Printf("server accept err: %v\n", err)
			return
		}

		//plugin aop
		conn, ok := l.Plugins.ConnAcceptHook(conn)
		if !ok {
			conn.Close()
			continue
		}
		log.Printf("server accepted conn: %v\n", conn.RemoteAddr().String())

		//create new routine worker each connection
		go l.handleConn(conn)
	}
}

//handle each connection
//TODO:对异常 err 处理
func (l *RPCListener) handleConn(conn net.Conn) {
	if l.isShutdown() {
		return
	}

	//catch panic
	defer func() {
		if err := recover(); err != nil {

			log.Printf("server %s catch panic err:%s\n", conn.RemoteAddr(), err)

			//closeConn(conn)
		}

		l.CloseConn(conn)
	}()

	for {
		if l.isShutdown() {
			return
		}

		//readtimeout
		startTime := time.Now()
		if l.option.ReadTimeout != 0 {
			conn.SetReadDeadline(startTime.Add(l.option.ReadTimeout))
		}

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
			log.Println("server decode request err:%v\n", err)
			return
		}
		log.Printf("server decode data finish:%v\n", inArgs)

		//handling
		atomic.AddInt32(&l.handlingNum, 1)
		defer atomic.AddInt32(&l.handlingNum, -1)

		//call local service
		handler, ok := l.Handlers[msg.ServiceClass]
		if !ok {
			log.Println("server can not found handler")
			return
		}

		l.Plugins.BeforeCallHook(msg.ServiceClass, msg.ServiceMethod, inArgs) //ctx
		result, err := handler.Handle(msg.ServiceMethod, inArgs)
		l.Plugins.AfterCallHook(msg.ServiceClass, msg.ServiceMethod, inArgs, result, err)
		log.Println("server call local service finish! result:", result)

		//encode
		encodeRes, err := coder.Encode(result) //[]byte result + err
		if err != nil {
			log.Printf("server encode err:%v\n", err)
			return
		}

		//send result
		if l.option.WriteTimeout != 0 {
			conn.SetWriteDeadline(startTime.Add(l.option.WriteTimeout))
		}

		l.Plugins.BeforeWriteHook(encodeRes)
		err = l.sendData(conn, encodeRes)
		l.Plugins.AfterWriteHook(encodeRes, err)
		if err != nil {
			log.Printf("server send err:%v\n", err) //timeout
			return
		}
		log.Printf("server send result finish! total runtime: %v", time.Now().Sub(startTime).Seconds())
	}
}

func (l *RPCListener) receiveData(conn net.Conn) (*protocol.RPCMsg, error) {
	l.Plugins.BeforeReadHook() //ctx

	msg, err := protocol.Read(conn)
	if err == io.EOF { //close
		log.Printf("server read finish:%v\n", err)
		return msg, nil
	}

	l.Plugins.AfterReadHook(msg, err)

	if err != nil {
		//rate limit
		return nil, err
	}

	//decode
	/*coder := global.Codecs[msg.Header.SerializeType()] //get from cache
	if coder == nil {
		return
	}
	inArgs := make([]interface{}, 0)
	err = coder.Decode(msg.Payload, &inArgs) //rpcdata
	if err != nil {
		log.Println("server decode request err:%v\n", err)
		return
	}
	log.Printf("server decode data finish:%v\n", inArgs)
	*/

	return msg, nil
}

func (l *RPCListener) sendData(conn net.Conn, payload []byte) error {
	resMsg := protocol.NewRPCMsg()
	resMsg.SetVersion(config.Protocol_MsgVersion)
	resMsg.SetMsgType(protocol.Response)
	resMsg.SetCompressType(protocol.None)
	resMsg.SetSerializeType(protocol.Gob)
	resMsg.Payload = payload
	return resMsg.Send(conn)
}

//net addr
func (l *RPCListener) GetAddrs() []string {
	//l.nl.Addr()
	addr := fmt.Sprintf("tcp://%s:%d", l.ServiceIp, l.ServicePort)
	return []string{addr}
}

func (l *RPCListener) getDoneChan() <-chan struct{} {
	return l.doneChan
}

func (l *RPCListener) closeDoneChan() {
	select {
	case <-l.doneChan:
	default:
		close(l.doneChan)
	}
}

func (l *RPCListener) CloseConn(conn net.Conn) {
	//activeconn

	conn.Close()

	//plugin
	log.Println("server closed")
}

func (l *RPCListener) Close() {
	if l.nl != nil {
		l.nl.Close()
	}
	l.closeDoneChan()
}

func (l *RPCListener) Shutdown() {
	atomic.CompareAndSwapInt32(&l.shutdown, 0, 1)
	for {
		if atomic.LoadInt32(&l.handlingNum) == 0 {
			break
		}
	}
	l.closeDoneChan()
	log.Println("server shutdown")
}

func (l *RPCListener) isShutdown() bool {
	return atomic.LoadInt32(&l.shutdown) == 1
}
