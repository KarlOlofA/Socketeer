package network

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/text/encoding/unicode"
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

func (p *Packet) BuildFromByteSlice(packetByteSlice []byte) error {
	reader := bytes.NewReader(packetByteSlice)

	fmt.Println("Go Input (hex):")
	fmt.Println(hex.EncodeToString(packetByteSlice))

	err := binary.Read(reader, binary.BigEndian, &p.Key)
	if err != nil {
		return err
	}

	var length uint32
	err = binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return err
	}

	err = binary.Read(reader, binary.BigEndian, &p.User)
	if err != nil {
		return err
	}

	data := make([]byte, length)
	_, err = reader.Read(data)
	if err != nil {
		return err
	}

	decoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()
	decodedKey, err := decoder.Bytes(data)
	if err != nil {
		return errors.New("Failed to decode key")
	}
	p.Key = string(decodedKey)

	remaining := reader.Len()
	p.Data = make([]byte, remaining)
	_, err = reader.Read(p.Data)

	return nil
}

func (p *Packet) FromByteSlice(slice []byte) error {

	if len(slice) < 25 {
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
