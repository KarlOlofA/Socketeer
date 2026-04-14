package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"golang.org/x/text/encoding/unicode"
)

type Packet struct {
	PacketType uint16 `json:"packetType"`
	Key        string `json:"key"`
	Data       []byte `json:"data"`
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

	err := binary.Read(reader, binary.BigEndian, &p.PacketType)
	if err != nil {
		return err
	}

	var keySize uint16
	err = binary.Read(reader, binary.BigEndian, &keySize)
	if err != nil {
		return err
	}

	keyBytes := make([]byte, keySize)
	_, err = reader.Read(keyBytes)
	if err != nil {
		return err
	}

	decoder := unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder()
	decodedKey, err := decoder.Bytes(keyBytes)
	if err != nil {
		return errors.New("Failed to decode key")
	}
	p.Key = string(decodedKey)

	remaining := reader.Len()
	p.Data = make([]byte, remaining)
	_, err = reader.Read(p.Data)

	return nil
}
