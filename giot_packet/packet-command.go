package giot_packet

import (
	"bytes"
)

type PacketCommand struct {
	Command Command
	AckId   uint8
	Payload *bytes.Buffer
}

func (packResp *PacketCommand) Encode(buffer *bytes.Buffer) error {
	pack := NewPacket()
	pack.PacketType = PACKET_TYPE_RESPONSE

	if err := EncodeData(&(packResp.Command), DT_BUFFER, pack.Payload); err != nil {
		return err
	}

	if err := EncodeData(&(packResp.AckId), DT_BYTE, pack.Payload); err != nil {
		return err
	}

	if packResp.Payload != nil {
		if err := EncodeData(packResp.Payload, DT_BUFFER, pack.Payload); err != nil {
			return err
		}
	}

	pack.Length = uint16(pack.Payload.Len())
	pack.Encode(buffer)
	return nil
}

func PacketCommandDecode(pack *Packet) *PacketCommand {
	packConn := &PacketCommand{}
	packConn.Payload = &bytes.Buffer{}

	buffer := bytes.NewReader(pack.Payload.Bytes())

	if _, err := DecodeData(&(packConn.Command), DT_BUFFER, buffer); err != nil {
		return nil
	}

	if _, err := DecodeData(&(packConn.AckId), DT_BYTE, buffer); err != nil {
		return nil
	}

	if _, err := DecodeData(packConn.Payload, DT_BUFFER, buffer); err != nil {
		return nil
	}

	return packConn
}
