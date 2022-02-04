package giotgo

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS:  -L./g-IoT/build/lib -lg-iot
// #include "./g-IoT/test/conf/config.h"
// #include "./g-IoT/src/include/g-iot/packet.h"
// #include <stdio.h>
// #include <stdlib.h>
/*

int GIoT_Packet_sz = sizeof(GIoT_Packet_t);

*/
import "C"
import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

type ptrAddr uint64

type packet struct {
	packetType byte
	length     uint16
	payload    *bytes.Buffer
}

func NewPacket() *packet {
	pck := &packet{}
	pck.payload = &bytes.Buffer{}
	return pck
}
func (pck *packet) reset() {
	pck.packetType = 0
	pck.length = 0
	pck.payload.Reset()
}

func (pck *packet) packetDecode(buffer *bytes.Buffer) (readlen int, err error) {
	readlen = 0
	var flag byte = 0
	for flag != byte(C.GIOT_PACKET_FLAG) {
		if buffer.Len() == 0 {
			return 0, &errPacketInvalid{}
		}

		flag, err = buffer.ReadByte()
		if err != nil {
			return
		}
	}

	buffer.UnreadByte()

	cPacket := C.GIoT_Packet_t{}
	bufferPtr := C.CBytes(buffer.Bytes())
	bufferLen := buffer.Len()

	tmpPacket := struct {
		PacketType byte
		Length     uint16
		PayloadPtr uint64
	}{}

	pck.reset()
	readlen = int(C.GIoT_Packet_Decode(&cPacket, (*C.uchar)(bufferPtr), C.ushort(bufferLen)))
	if readlen == 0 {
		err = &errPacketInvalid{}
		return
	}
	cPacketBytes := C.GoBytes(unsafe.Pointer(&cPacket), C.GIoT_Packet_sz)
	if err = binary.Read(bytes.NewReader(cPacketBytes), binary.LittleEndian, &tmpPacket); err != nil {
		return
	}

	pck.packetType = tmpPacket.PacketType
	pck.length = tmpPacket.Length

	bufferOffset := readlen
	if tmpPacket.PayloadPtr != 0 {
		bufferOffset = int(uintptr(tmpPacket.PayloadPtr) - uintptr(bufferPtr))
	}
	pck.payload.Write(buffer.Bytes()[bufferOffset:])
	buffer.Next(bufferOffset)

	C.free(bufferPtr)
	return
}
