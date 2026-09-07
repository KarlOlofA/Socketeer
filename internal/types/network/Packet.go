package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Packet struct {
	Key    string
	Length uint32
	User   string
	Data   []byte
}

const (
	Data    uint16 = 0
	Welcome uint16 = 1
)

func (p Packet) ToByteSlice() ([]byte, error) {

	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, p)
	if err != nil {
		return nil, fmt.Errorf("Packet parse failed -> %v\n", err)
	}

	return buffer.Bytes(), nil

}

func (p *Packet) FromByteSlice(slice []byte) error {

	if slice == nil {
		return fmt.Errorf("byte slice is nil")
	}

	if len(slice) < 25 {
		fmt.Printf("%v\n", string(slice))
		return fmt.Errorf("Byte slice to small")
	}

	p.Key = string(slice[:4])
	length := slice[4:8]
	var num uint32
	err := binary.Read(bytes.NewReader(length), binary.BigEndian, &num)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if num > 1000 {
		return fmt.Errorf("Message length to large: %d", num)
	}

	p.Length = num
	p.User = string(slice[8:24])
	p.Data = slice[24 : 24+num]
	return nil
}
