package giot_packet

import "bytes"

type PacketConnack struct {
	Status     RespStatus
	properties [MAX_PROP_CONNACK]Property
}

func NewPacketConnack(respStatus RespStatus) *PacketConnack {
	packConnack := &PacketConnack{
		Status: respStatus,
	}
	Properties(packConnack.properties[:]).AddProperty(PROP_PROTOCOL_V, byte(PROTOCOL_VERSION))

	return packConnack
}

func (packConnack *PacketConnack) Encode(buffer *bytes.Buffer) error {
	pack := NewPacket(buffer.Bytes()[16:])
	pack.PacketType = PACKET_TYPE_CONNACK

	if err := EncodeData(&(packConnack.Status), DT_BYTE, pack.Payload); err != nil {
		return err
	}
	if err := Properties(packConnack.properties[:]).Encode(pack.Payload); err != nil {
		return err
	}
	pack.Length = uint16(pack.Payload.Len())
	pack.Encode(buffer)
	return nil
}
