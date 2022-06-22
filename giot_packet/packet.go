package giot_packet

import (
	"bytes"
	"encoding/binary"
)

type PacketHeader struct {
	FlagStart  uint8
	PacketType PacketType
	FlagClose  [2]uint8
}

type Packet struct {
	PacketType PacketType
	Length     uint16
	Payload    *bytes.Buffer
}

func NewPacket(payload []byte) *Packet {
	pack := &Packet{}
	pack.Payload = bytes.NewBuffer(payload)
	pack.Payload.Reset()
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
		return errInvalidPacket
	}

	buffer.Reset()

	header := PacketHeader{
		PACKET_FLAG_START,
		pack.PacketType,
		PACKET_FLAG_CLOSE,
	}

	bufSpaceLen := buffer.Cap() - buffer.Len()
	if bufSpaceLen < binary.Size(header) {
		return errBufferNoSpace
	}

	if err := binary.Write(buffer, binary.BigEndian, header); err != nil {
		return err
	}

	if err := EncodeData(&pack.Length, DT_DINT, buffer); err != nil {
		return err
	}

	bufSpaceLen = buffer.Cap() - buffer.Len()
	if bufSpaceLen < pack.Payload.Len() {
		return errBufferNoSpace
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
	if header.FlagStart != PACKET_FLAG_START {
		return errInvalidPacket
	}
	if header.FlagClose != PACKET_FLAG_CLOSE {
		return errInvalidPacket
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
