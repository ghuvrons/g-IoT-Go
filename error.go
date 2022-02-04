package giotgo

type errPacketInvalid struct{}

func (err *errPacketInvalid) Error() string {
	return "Packet Invalid"
}
