package giot_packet

import (
	"bytes"
	"reflect"
)

type Property struct {
	id   byte
	data Data
}

type Properties []Property

func (properties Properties) Encode(buf *bytes.Buffer) error {
	tmpBuffer := &bytes.Buffer{}
	for i := range properties {
		if properties[i].id == 0 {
			break
		}

		// write ID
		if err := tmpBuffer.WriteByte(properties[i].id); err != nil {
			return err
		}

		// write data
		data := properties[i].data
		switch data.(type) {
		case []byte:
			data = bytes.NewReader(data.([]byte))
			tmpBuffer.Grow(len(data.([]byte)) + 4)
		case string:
			tmpBuffer.Grow(len(data.(string)) + 4)
		case *bytes.Reader:
			tmpBuffer.Grow(data.(*bytes.Reader).Len() + 4)
		case *bytes.Buffer:
			tmpBuffer.Grow(data.(*bytes.Buffer).Len() + 4)
		default:
			tmpBuffer.Grow(8)
		}
		if err := EncodeData(data, (&properties[i]).getDataType(), tmpBuffer); err != nil {
			return err
		}
	}

	if err := EncodeData(tmpBuffer.Len(), DT_DINT, buf); err != nil {
		return err
	}

	bufSpaceLen := buf.Cap() - buf.Len()
	if bufSpaceLen < tmpBuffer.Len() {
		return errBufferNoSpace
	}
	if _, err := tmpBuffer.WriteTo(buf); err != nil {
		return err
	}

	return nil
}

func (properties Properties) Decode(buffer *bytes.Reader) (readLen int, err error) {
	var propBufLen int
	var tmpReadLen int
	readLen = 0
	err = nil

	tmpReadLen, err = DecodeData(&propBufLen, DT_DINT, buffer)
	readLen += tmpReadLen
	if err != nil {
		return
	}

	for i := range properties {
		if propBufLen <= 0 {
			break
		}
		if buffer.Len() <= 0 {
			return readLen, errInvalidData
		}
		// read id
		properties[i].id, err = buffer.ReadByte()
		readLen += 1
		propBufLen -= 1
		if err != nil {
			return
		}

		// read data
		if properties[i].data == nil {
			switch (&properties[i]).getDataType() {
			case DT_BYTE:
				properties[i].data = byte(0)
			case DT_DINT:
				properties[i].data = int(0)
			case DT_UTF8, DT_BINARY:
				properties[i].data = &bytes.Buffer{}
			}
		}
		var data Data = properties[i].data
		if reflect.ValueOf(data).Kind() != reflect.Ptr {
			data = &data
		}
		tmpReadLen, err = DecodeData(data, (&properties[i]).getDataType(), buffer)
		readLen += tmpReadLen
		propBufLen -= tmpReadLen
		if err != nil {
			return
		}

	}
	return
}

func (prop *Property) getDataType() dataType {
	var r dataType = 0
	switch prop.id {
	case PROP_AUTH_USER, PROP_AUTH_PASS:
		r = DT_UTF8
	}
	return r
}

func (properties Properties) AddProperty(id uint8, data interface{}) {
	for i := range properties {
		if properties[i].id != 0 {
			continue
		}

		switch data.(type) {
		case string:
			data = []byte(data.(string))
		}

		properties[i].id = id
		properties[i].data = data
		break
	}
}
