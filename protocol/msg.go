package protocol

import (
	"encoding/binary"
	"fmt"
	"io"
	"unsafe"
)

const (
	SPLIT_LEN = 4
)

type RPCMsg struct {
	*Header
	ServiceMethod string
	Payload       []byte
	Metadata      map[string]string
}

func NewRPCMsg() *RPCMsg {
	header := Header([HEADER_LEN]byte{})
	header[0] = magicNumber
	return &RPCMsg{
		Header: &header,
	}
}

func (msg *RPCMsg) Send(writer io.Writer) (int64, error) {
	//send header
	h, err := writer.Write(msg.Header[:])
	if err != nil {
		return int64(h), err
	}

	//write body total len
	dataLen := SPLIT_LEN + len(msg.ServiceMethod) + SPLIT_LEN + len(msg.Payload)
	err = binary.Write(writer, binary.BigEndian, uint32(dataLen))
	if err != nil {
		return 0, err
	}

	//write service len 4
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceMethod)))
	if err != nil {
		return 0, err
	}

	//write service content
	err = binary.Write(writer, binary.BigEndian, StringToByte(msg.ServiceMethod))
	if err != nil {
		return 0, err
	}

	//write payload len 4
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.Payload)))
	if err != nil {
		return 0, err
	}

	//write payload
	//err = binary.Write(writer, binary.BigEndian, msg.Payload)
	nc, err := writer.Write(msg.Payload)
	if err != nil {
		return 0, err
	}
	return int64(nc), nil
}

func (msg *RPCMsg) Decode(r io.Reader) error {
	_, err := io.ReadFull(r, msg.Header[:]) //magicNumber
	if !msg.Header.CheckMagicNumber() {
		return fmt.Errorf("magic number error: %v", msg.Header[0])
	}

	headerByte := make([]byte, 4)
	_, err = io.ReadFull(r, headerByte)
	if err != nil {
		return err
	}

	//datalen
	dataLen := binary.BigEndian.Uint32(headerByte)

	//all data byte
	data := make([]byte, dataLen)
	_, err = io.ReadFull(r, data)

	//4
	n := 0
	l := binary.BigEndian.Uint32(data[n : n+4]) //0,4

	//serviceMethod
	n = n + 4 //
	nEnd := n + int(l)
	msg.ServiceMethod = ByteToString(data[n:nEnd])

	//4
	n = nEnd
	l = binary.BigEndian.Uint32(data[n : n+4])

	//payload
	n = n + 4
	msg.Payload = data[n:]

	return nil
}

func Read(r io.Reader) (*RPCMsg, error) {
	msg := NewRPCMsg()
	err := msg.Decode(r)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func StringToByte(s string) []byte {
	r := (*[2]uintptr)(unsafe.Pointer(&s))
	k := [3]uintptr{r[0], r[1], r[1]}
	return *(*[]byte)(unsafe.Pointer(&k))
}

func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
