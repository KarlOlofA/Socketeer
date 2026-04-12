package types

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Packet struct {
	KeySize    uint16
	Key        string `json:"key"`
	PacketType int    `json:"packetType"`
	Data       []byte `json:"data"`
}

const (
	Data    int = 0
	Welcome int = 1
)

type TestPacket struct {
	KeySize    uint16
	Key        string
	PacketType int
	Data       []byte
}

func (p Packet) ToByteSlice() ([]byte, error) {

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(p)
	if err != nil {
		return nil, fmt.Errorf("Packet parse failed -> %v\n", err)
	}

	return buffer.Bytes(), nil

}

func (p *Packet) BuildFromByteSlice(packetByteSlice []byte) error {
	buffer := bytes.NewBuffer(packetByteSlice)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&p)
	if err != nil {
		return err
	}
	return nil
}
