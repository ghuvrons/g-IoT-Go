package giot_packet

type errDataInvalid struct{}

func (err *errDataInvalid) Error() string {
	return "Data Invalid"
}

type errPacketInvalid struct{}

func (err *errPacketInvalid) Error() string {
	return "Packet Invalid"
}

type errAuthInvalid struct{}

func (err *errAuthInvalid) Error() string {
	return "username or password error"
}
