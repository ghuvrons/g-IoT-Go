package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"os"

	giotgo "github.com/ghuvrons/g-IoT-Go"
	giot_packet "github.com/ghuvrons/g-IoT-Go/giot_packet"
)

func setCmdHandlers(server *giotgo.Server) {
	server.On(CMD_GET_INFO,
		func(client *giotgo.ClientHandler, data giot_packet.Data) (giot_packet.RespStatus, *bytes.Buffer) {
			b, isOK := data.(*bytes.Buffer)
			if !isOK {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			length, crc := calcInfoFile(b.String())

			fmt.Println(length, crc)

			buf := &bytes.Buffer{}
			binary.Write(buf, binary.BigEndian, length)
			binary.Write(buf, binary.BigEndian, crc)

			return giot_packet.RESP_OK, buf
		},
	)

	readFileBuffer := make([]byte, 1024)
	var crc uint32 = 0
	server.On(CMD_DOWNLOAD,
		func(client *giotgo.ClientHandler, data giot_packet.Data) (giot_packet.RespStatus, *bytes.Buffer) {
			b, isOK := data.(*bytes.Buffer)
			if !isOK {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			offset := binary.BigEndian.Uint32(b.Bytes()[:4])
			b.Next(4)
			readLen := binary.BigEndian.Uint32(b.Bytes()[:4])
			b.Next(4)
			path := b.String()

			f, err := os.Open(path)
			defer func() {
				f.Close()
			}()
			if err != nil {
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			if ret, err := f.Seek(int64(offset), 0); err != nil {
				fmt.Println(err)
				fmt.Println(offset, "=>", ret)
			} else {
			}

			if readLen > uint32(cap(readFileBuffer)) {
				readLen = uint32(cap(readFileBuffer))
			}
			n2, err := f.Read(readFileBuffer)

			// debugging
			crc = crc32.Update(crc, crc32.IEEETable, readFileBuffer[:n2])
			fmt.Printf("%d|0x%.2X\r\n", offset, crc)
			if offset == 413696 {
				// fmt.Print("hayooo", offset)
			}

			if err != nil {
				fmt.Println(err)
				return giot_packet.RESP_UNKNOWN_ERROR, nil
			}

			buf := bytes.NewBuffer(readFileBuffer[:n2])

			return giot_packet.RESP_OK, buf
		},
	)
}

// Calculate length and crc of file
func calcInfoFile(path string) (uint32, uint32) {
	var length uint32 = 0
	var crc uint32 = 0

	f, err := os.Open(path)

	defer func() {
		f.Close()
	}()

	if err != nil {
		return 0, 0
	}

	b2 := make([]byte, 256)

	for true {
		n2, err := f.Read(b2)
		if err != nil {
			if err == io.EOF {
				break
			}
			return 0, 0
		}

		length += uint32(n2)
		if n2 == 0 {
			break
		}
		crc = crc32.Update(crc, crc32.IEEETable, b2[:n2])
		if n2 < 256 {
			break
		}
	}
	return length, crc
}
