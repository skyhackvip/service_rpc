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
	ServiceClass  string
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

	//write body total len :4 byte
	dataLen := SPLIT_LEN + len(msg.ServiceClass) + SPLIT_LEN + len(msg.ServiceMethod) + SPLIT_LEN + len(msg.Payload)
	err = binary.Write(writer, binary.BigEndian, uint32(dataLen)) //4
	if err != nil {
		return 0, err
	}

	//write service class len :4 byte
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceClass)))
	if err != nil {
		return 0, err
	}

	//write service class
	err = binary.Write(writer, binary.BigEndian, StringToByte(msg.ServiceClass))
	if err != nil {
		return 0, err
	}

	//write service method len :4 byte
	err = binary.Write(writer, binary.BigEndian, uint32(len(msg.ServiceMethod)))
	if err != nil {
		return 0, err
	}

	//write service method
	err = binary.Write(writer, binary.BigEndian, StringToByte(msg.ServiceMethod))
	if err != nil {
		return 0, err
	}

	//write payload len :4 byte
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
	//read header
	_, err := io.ReadFull(r, msg.Header[:])
	if !msg.Header.CheckMagicNumber() { //magicNumber
		return fmt.Errorf("magic number error: %v", msg.Header[0])
	}

	//total body len
	headerByte := make([]byte, 4)
	_, err = io.ReadFull(r, headerByte)
	if err != nil {
		return err
	}
	bodyLen := binary.BigEndian.Uint32(headerByte)

	//read all body
	data := make([]byte, bodyLen)
	_, err = io.ReadFull(r, data)

	//service class len
	start := 0
	end := start + SPLIT_LEN
	classLen := binary.BigEndian.Uint32(data[start:end]) //0,4

	//service class
	start = end
	end = start + int(classLen)
	msg.ServiceClass = ByteToString(data[start:end]) //4,x

	//service method len
	start = end
	end = start + SPLIT_LEN
	methodLen := binary.BigEndian.Uint32(data[start:end]) //x,x+4

	//service method
	start = end
	end = start + int(methodLen)
	msg.ServiceMethod = ByteToString(data[start:end]) //x+4, x+4+y

	//payload len
	start = end
	end = start + SPLIT_LEN
	binary.BigEndian.Uint32(data[start:end]) //x+4+y, x+y+8 payloadLen

	//payload
	start = end
	msg.Payload = data[start:]
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
