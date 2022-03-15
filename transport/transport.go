package transport

import (
	"encoding/binary"
	"github.com/skyhackvip/service_rpc/config"
	"io"
	"net"
)

type Transport struct {
	conn net.Conn
}

func NewTransport(conn net.Conn) *Transport {
	return &Transport{conn}
}

//read data
func (trans *Transport) Read() ([]byte, error) {
	//header
	header := make([]byte, config.HEADER_LEN) //4
	_, err := io.ReadFull(trans.conn, header)
	if err != nil {
		return nil, err
	}

	//data
	dataLen := binary.BigEndian.Uint32(header) //header保存的数据为data的长度
	data := make([]byte, dataLen)
	_, err = io.ReadFull(trans.conn, data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//send data
func (trans *Transport) Send(data []byte) error {
	headerLen := config.HEADER_LEN
	buffer := make([]byte, headerLen+len(data))
	//长度
	binary.BigEndian.PutUint32(buffer[:headerLen], uint32(len(data)))

	copy(buffer[headerLen:], data)
	_, err := trans.conn.Write(buffer)
	return err
}
