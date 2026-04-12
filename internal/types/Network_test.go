package types

import (
	"encoding/binary"
	"fmt"
	"testing"
)

func TestPacketToByteSlice(t *testing.T) {

	size := []byte("1234")
	keySize := binary.LittleEndian.Uint16(size)
	packet := Packet{
		KeySize:    keySize,
		Key:        "1234",
		PacketType: 1,
		Data:       []byte("Test"),
	}

	slice, err := packet.ToByteSlice()
	if err != nil {
		t.Errorf("%v\n", err)
	}

	fmt.Printf("Byte Slice -> %v\n", slice)

	fmt.Printf("Slice Data -> %s\n", string(slice))
}

func TestByteSliceToPacket(t *testing.T) {
	size := []byte("1234")
	keySize := binary.LittleEndian.Uint16(size)

	packet := Packet{
		KeySize:    keySize,
		Key:        "1234",
		PacketType: 1,
		Data:       []byte("Test"),
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
	fmt.Printf("Key -> %s\n PacketType -> %d\n Data -> %s\n", packet.Key, packet.PacketType, string(packet.Data))
}
