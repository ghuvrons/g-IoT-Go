package giot_packet

import (
	"bytes"
	"encoding/binary"
)

type PacketHeader struct {
	Flag       uint8
	PacketType PacketType
}

type Packet struct {
	PacketType PacketType
	Length     uint16
	Payload    *bytes.Buffer
}

func NewPacket() *Packet {
	pack := &Packet{}
	pack.Payload = &bytes.Buffer{}
	return pack
}

func (pack *Packet) Reset() {
	pack.PacketType = 0
	pack.Length = 0
	pack.Payload.Reset()
}

func (pack *Packet) IsValid() bool {
	if pack.PacketType == 0 {
		return false
	}

	if int(pack.Length) != pack.Payload.Len() {
		return false
	}

	return true
}

func (pack *Packet) Encode(buffer *bytes.Buffer) error {
	if !pack.IsValid() {
		return &errPacketInvalid{}
	}

	header := PacketHeader{
		PACKET_FLAG,
		pack.PacketType,
	}

	buffer.Reset()
	if err := binary.Write(buffer, binary.BigEndian, header); err != nil {
		return err
	}

	if err := EncodeData(&pack.Length, DT_DINT, buffer); err != nil {
		return err
	}

	if _, err := pack.Payload.WriteTo(buffer); err != nil {
		return err
	}

	return nil
}

func (pack *Packet) Decode(buffer *bytes.Buffer) error {
	var readLen int = 0
	var tmpReadLen int = 0

	// read header
	header := PacketHeader{}
	if err := binary.Read(bytes.NewReader(buffer.Bytes()), binary.BigEndian, &header); err != nil {
		return err
	}
	if header.Flag != PACKET_FLAG {
		return &errPacketInvalid{}
	}

	tmpReadLen = binary.Size(&header)
	buffer.Next(tmpReadLen)
	readLen += tmpReadLen

	pack.PacketType = header.PacketType

	// read Length
	var err error
	tmpReadLen, err = DecodeData(&pack.Length, DT_DINT, bytes.NewReader(buffer.Bytes()))

	if err != nil {
		// unread header if error
		for i := 0; i < readLen; i++ {
			buffer.UnreadByte()
		}
		return err
	}
	buffer.Next(tmpReadLen)
	readLen += tmpReadLen

	// read and write to payload
	pack.Payload.Reset()
	if buffer.Cap() < int(pack.Length) {
		tmpReadLen = buffer.Cap()
	} else {
		tmpReadLen = int(pack.Length)
	}
	pack.Payload.Write(buffer.Bytes()[:tmpReadLen])
	buffer.Next(tmpReadLen)
	readLen += tmpReadLen

	return nil
}
