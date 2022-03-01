package data

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type RPCDataIf interface {
	Encode(data RPCData) ([]byte, error)
	Decode(data []byte) (RPCData, error)
}

type RPCData struct {
	Name string
	Args []interface{}
	Err  string
}

type GobData struct {
}

func (g *GobData) Encode(data RPCData) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(data); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (g *GobData) Decode(data []byte) (RPCData, error) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)

	var resData RPCData
	if err := decoder.Decode(&resData); err != nil {
		return RPCData{}, nil
	}
	return resData, nil
}

type JsonData struct {
}

func (j *JsonData) Encode(data RPCData) ([]byte, error) {
	return json.Marshal(data)
}

func (j *JsonData) Decode(data []byte) (RPCData, error) {
	r := RPCData{}
	err := json.Unmarshal(data, &r)
	return r, err
}

func New(kind string) RPCDataIf {
	switch kind {
	case "gob":
		return &GobData{}
	case "json":
		return &JsonData{}
	default:
		return &JsonData{}
	}
}
