package protocol

import (
	"encoding/binary"
	"errors"
)

const (
    PACKET_HEADER_SIZE = 4
    MAX_PACKET_LENGTH = 4096
)

var ErrorMaxSizeExceeded = errors.New("Max packet Size exceeded")
var ErrorEmptyPacket = errors.New("Packet size 0. Packet is Empty")

type PayloadType uint8

const (
	MSG_STATUS_READ PayloadType = 1
	MSG_STATUS_DELV
	MSG_STATUS_SENT
    MSG_SEND
    MSG_RECEIVE
	CONN_ERR
	CONN_FIN
	CONN_INIT
)

type Packet struct {
	Version       uint8
	PayloadType   PayloadType
	PayloadLength uint16
	Payload       []byte
}

func (p *Packet) ToBytes() []byte {
    output := make([]byte, PACKET_HEADER_SIZE + len(p.Payload))
    output[0] = p.Version
    output[1] = uint8(p.PayloadType)
    binary.BigEndian.PutUint16(output[2:], p.PayloadLength)
    copy(output[4:], p.Payload)

    return output
}

func NewPacket(payloadType PayloadType, payload Payload) Packet {
    payloadBinary, _ := payload.MarshalBinary()
	packet := Packet{
        Version: 1,
        PayloadType: payloadType,
        PayloadLength: uint16(payload.Length()),
        Payload: payloadBinary,
	}

	return packet
}

func PacketFromBytes(data []byte) *Packet {
	output := &Packet{}

	output.Version = data[0]
	output.PayloadType = PayloadType(data[1])
	output.PayloadLength = binary.BigEndian.Uint16(data[2:])
	output.Payload = make([]byte, output.PayloadLength)
	copy(output.Payload, data[4:])

	return output
}
