package giotgo

// #cgo CFLAGS: -g -Wall
// #cgo LDFLAGS:  -L./g-IoT/build/lib -lg-iot
// #include "./g-IoT/test/conf/config.h"
// #include "./g-IoT/src/include/g-iot/packet.h"
// #include <stdio.h>
/*
void gUpdate(void) {
	printf("sizeof: %d\r\n", (int) sizeof(GIoT_Packet_t));
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type GIoT_Packet struct {
	packetType byte
	length     uint16
	payload    []byte
}

func PacketEncode() {
	C.gUpdate()
	packet := &GIoT_Packet{}
	buffer := []byte{0xC7, 0x01, 5, 13, 42, 12, 33, 21}
	packetPtr := (*C.GIoT_Packet_t)(unsafe.Pointer(packet))
	bufferPtr := (*C.uchar)(C.CBytes(buffer))
	C.GIoT_Packet_Decode(packetPtr, bufferPtr, 8)
	packet2 := C.GoBytes(unsafe.Pointer(packetPtr), 12)
	buffer2 := C.GoBytes(unsafe.Pointer(bufferPtr), 12)
	fmt.Println(packet2, buffer2, packet)
}
