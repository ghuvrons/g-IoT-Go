package giot_packet

import (
	"bytes"
)

type PacketResponse struct {
	Status  RespStatus
	AckId   uint8
	Payload *bytes.Buffer
}

func NewPacketResponse(respStatus RespStatus) *PacketResponse {
	packResp := &PacketResponse{}
	packResp.Status = respStatus
	return packResp
}

func (packResp *PacketResponse) Encode(buffer *bytes.Buffer) error {
	pack := NewPacket()
	pack.PacketType = PACKET_TYPE_RESPONSE

	if err := EncodeData(&(packResp.Status), DT_BYTE, pack.Payload); err != nil {
		return err
	}

	if err := EncodeData(&(packResp.AckId), DT_BYTE, pack.Payload); err != nil {
		return err
	}

	if packResp.Payload != nil {
		if err := EncodeData(bytes.NewReader(packResp.Payload.Bytes()), DT_BUFFER, pack.Payload); err != nil {
			return err
		}
	}

	pack.Length = uint16(pack.Payload.Len())
	pack.Encode(buffer)
	return nil
}

func PacketResponseDecode(pack *Packet) *PacketResponse {
	packConn := &PacketResponse{}
	packConn.Payload = &bytes.Buffer{}

	buffer := bytes.NewReader(pack.Payload.Bytes())

	if _, err := DecodeData(&(packConn.Status), DT_BYTE, buffer); err != nil {
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
