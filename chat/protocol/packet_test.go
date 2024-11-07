package protocol

import (
	"testing"
	"time"
)

func TestNewPacketMessageReceive(t *testing.T) {
	input := &MessageReceieve{
		Id:          1,
		SenderId:    2,
		RecipientId: 3,
		Timestamp:   time.Unix(12345678, 0),
		Content:     "hello",
	}

	result := NewPacket(MSG_RECEIVE, input)

	expectedPayload := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
		0x68, 0x65, 0x6C, 0x6C, 0x6F,
	}

	expected := Packet{
		Version:       1,
		PayloadType:   MSG_RECEIVE,
		PayloadLength: 37,
		Payload:       expectedPayload,
	}

	if result.Version != expected.Version {
		t.Errorf("Expected version: %d. Got: %d", expected.Version, result.Version)
	}

	if result.PayloadType != expected.PayloadType {
		t.Errorf("Expected payload type: %d. Got: %d", expected.PayloadType, result.PayloadType)
	}

	if result.PayloadLength != expected.PayloadLength {
		t.Errorf("Expected payload length: %d. Got: %d", expected.PayloadLength, result.PayloadLength)
	}

	testMessageReceivePayload(result.Payload, expected.Payload, t)
}

func TestNewPacketMessageStatusSent(t *testing.T) {
	input := &MessageStatusSent{
		MessageId: 1234,
		Timestamp: time.Unix(12345678, 0),
	}

	result := NewPacket(MSG_STATUS_SENT, input)

	expectedPayload := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
	}

	expected := Packet{
		Version:       1,
		PayloadType:   MSG_STATUS_SENT,
		PayloadLength: 16,
		Payload:       expectedPayload,
	}

	if result.Version != expected.Version {
		t.Errorf("Expected version: %d. Got: %d", expected.Version, result.Version)
	}

	if result.PayloadType != expected.PayloadType {
		t.Errorf("Expected payload type: %d. Got: %d", expected.PayloadType, result.PayloadType)
	}

	if result.PayloadLength != expected.PayloadLength {
		t.Errorf("Expected payload length: %d. Got: %d", expected.PayloadLength, result.PayloadLength)
	}

	testMessageStatusSentPayload(result.Payload, expected.Payload, t)
}

// 8 milk
// 7-8 corriander
// 15 Cucumber
