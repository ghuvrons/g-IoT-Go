package giot_packet

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

type Data interface{}

func EncodeData(data Data, dt dataType, buf *bytes.Buffer) error {
	var rv reflect.Value
	bufSpaceLen := buf.Cap() - buf.Len()

	if bufSpaceLen == 0 {
		return errBufferNoSpace
	}

	rv = reflect.ValueOf(data)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return errInvalidData
	}

	switch dt {
	case DT_BYTE:
		var b []byte
		switch rk := rv.Kind(); rk {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			b = []byte{byte(rv.Uint())}
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			b = []byte{byte(rv.Int())}
		default:
			return errInvalidData
		}
		buf.Write(b)

	case DT_DINT:
		x64 := uint64(0)
		switch rk := rv.Kind(); rk {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			x64 = rv.Uint()
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			x64 = uint64(rv.Int())
		default:
			return errInvalidData
		}
		var b []byte

		if x64 < 0xFE {
			b = []byte{byte(x64)}
		} else if x64 <= 0xFFFF {
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, x64)
			buf.Write([]byte{byte(0xFE)})
			b = b[6:]
		} else {
			b = make([]byte, 8)
			binary.BigEndian.PutUint64(b, x64)
			buf.Write([]byte{byte(0xFF)})
			b = b[4:]
		}
		bufSpaceLen = buf.Cap() - buf.Len()
		if bufSpaceLen < len(b) {
			return errBufferNoSpace
		}
		if _, err := buf.Write(b); err != nil {
			return err
		}

	case DT_BINARY, DT_UTF8:
		rdr, isOK := data.(*bytes.Reader)
		if !isOK {
			return errInvalidData
		}
		if err := EncodeData(rdr.Len(), DT_DINT, buf); err != nil {
			return err
		}
		bufSpaceLen = buf.Cap() - buf.Len()
		if bufSpaceLen < rdr.Len() {
			return errBufferNoSpace
		}
		if _, err := rdr.WriteTo(buf); err != nil {
			return err
		}

	case DT_BUFFER:
		switch data.(type) {
		case *bytes.Reader:
			rdr := data.(*bytes.Reader)
			if bufSpaceLen < rdr.Len() {
				return errBufferNoSpace
			}
			if _, err := rdr.WriteTo(buf); err != nil {
				return err
			}

		case []byte:
			dataBytes := data.([]byte)
			if bufSpaceLen < len(dataBytes) {
				return errBufferNoSpace
			}
			if _, err := buf.Write(dataBytes); err != nil {
				return err
			}

		default:
			dataBytes := make([]byte, 8)
			tmpBuffer := bytes.NewBuffer(dataBytes)
			tmpBuffer.Reset()
			if err := binary.Write(tmpBuffer, binary.BigEndian, data); err != nil {
				return err
			}
			if bufSpaceLen < tmpBuffer.Len() {
				return errBufferNoSpace
			}
			if _, err := tmpBuffer.WriteTo(buf); err != nil {
				return err
			}
		}
	}
	return nil
}

func DecodeData(data Data, dt dataType, rdr *bytes.Reader) (readLen int, err error) {
	var n int = 0
	var x8 byte
	var rv reflect.Value
	readLen = 0
	err = nil

	rv = reflect.ValueOf(data)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if !rv.IsValid() || !rv.CanSet() {
		return readLen, errInvalidData
	}

	switch dt {
	case DT_BYTE:
		if rv.Kind() != reflect.Uint8 {
			return readLen, errInvalidData
		}

		x8, err = rdr.ReadByte()
		readLen += 1
		if err != nil {
			return
		}

		x64 := uint64(x8)
		if !rv.OverflowUint(x64) {
			rv.SetUint(x64)
		}

	case DT_DINT:
		x8, err = rdr.ReadByte()
		readLen += 1
		if err != nil {
			return
		}
		x64 := int64(0)

		if x8 == 0xFE {
			x64, err = readInt(rdr, 2)
			readLen += 2
		} else if x8 == 0xFF {
			x64, err = readInt(rdr, 4)
			readLen += 4
		} else {
			x64 = int64(x8)
		}
		if err != nil {
			return
		}
		switch rk := rv.Kind(); rk {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			if !rv.OverflowUint(uint64(x64)) {
				rv.SetUint(uint64(x64))
			}

		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			if !rv.OverflowInt(x64) {
				rv.SetInt(x64)
			}
		}

	case DT_BINARY, DT_UTF8:
		dataLength := 0

		n, err = DecodeData(&dataLength, DT_DINT, rdr)
		readLen += n

		if err != nil {
			return
		}

		if rdr.Len() < dataLength {
			return readLen, errInvalidData
		}

		dw, isOK := data.(io.Writer)
		if !isOK {
			return readLen, errInvalidData
		}

		tmpBuf := make([]byte, dataLength)
		rdr.Read(tmpBuf)
		n, err = dw.Write(tmpBuf)
		readLen += n

		if err != nil {
			return
		}

	case DT_BUFFER:
		switch data.(type) {
		case io.Writer:
			dw, isOK := data.(io.Writer)
			if !isOK {
				return readLen, errInvalidData
			}
			var n64 int64
			n64, err = rdr.WriteTo(dw)
			readLen += int(n64)

			if err != nil {
				return
			}

		default:
			dataSize := binary.Size(data)
			if dataSize > rdr.Len() {
				return readLen, errInvalidData
			}
			err = binary.Read(rdr, binary.BigEndian, data)
			readLen += dataSize

			if err != nil {
				return
			}
		}
	}

	return readLen, nil
}

// readUint read len(length) data as uint64
func readUint(rdr io.Reader, len int) (uint64, error) {
	// sometimes, data recived in length less than 8
	b := make([]byte, len)
	err := binary.Read(rdr, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}
	var result uint64 = 0

	switch len {
	case 1:
		result = uint64(b[0])
	case 2:
		result = uint64(binary.LittleEndian.Uint16(b))
	case 3, 4:
		result = uint64(binary.LittleEndian.Uint32(b))
	default:
		result = binary.LittleEndian.Uint64(b)
	}

	return result, nil
}

// readInt read len(length) data as int64 (signed int)
func readInt(rdr io.Reader, len int) (int64, error) {
	// sometimes, data recived in length less than 8
	b := make([]byte, len)
	err := binary.Read(rdr, binary.BigEndian, &b)
	if err != nil {
		return 0, err
	}

	var result int64 = 0

	switch len {
	case 1:
		result = int64(int8(b[0]))
	case 2:
		result = int64(int16(binary.LittleEndian.Uint16(b)))
	case 3, 4:
		result = int64(int32(binary.LittleEndian.Uint32(b)))
	default:
		result = int64(binary.LittleEndian.Uint64(b))
	}

	return result, nil
}
