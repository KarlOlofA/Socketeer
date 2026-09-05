package network

import (
	"encoding/binary"
	"fmt"
	"testing"
)

const (
	Key    = 4
	Length = 4
	User   = 16
)

func TestPacketToByteSlice(t *testing.T) {
	packet := Packet{
		Key:    "1234",
		Length: 1,
		Data:   []byte("Test"),
	}

	slice, err := packet.ToByteSlice()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Printf("Byte Slice -> %v\n", slice)

	fmt.Printf("Slice Data -> %s\n", string(slice))
}

func TestByteSliceToPacket(t *testing.T) {
	size := []byte("Test")
	keySize := binary.BigEndian.Uint32(size)

	packet := Packet{
		Key:    "1234",
		Length: keySize,
		User:   "",
		Data:   []byte("Test"),
	}

	slice, err := packet.ToByteSlice()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	var newPacket Packet
	err = newPacket.BuildFromByteSlice(slice)
	if err != nil {
		t.Errorf("Failed to build from byte slice -> %v\n", err)
	}
	packet = newPacket
	fmt.Printf("Key -> %s\n length -> %d\nuser -> %s\n Data -> %s\n", packet.Key, packet.Length, packet.User, string(packet.Data))
}
