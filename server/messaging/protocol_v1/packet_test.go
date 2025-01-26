package protocol_v1

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

func TestNewPacketMessageSentSuccess(t *testing.T) {
	input := &MessageSentSuccess{
		MessageId: 1234,
        RecipientId: 4321,
		Timestamp: time.Unix(12345678, 0),
        SentAt: time.Unix(12345678, 0),
	}

	result := NewPacket(MSG_SENT_SUCCESS, input)

	expectedPayload := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0xd2,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0xe1,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
	}

	expected := Packet{
		Version:       1,
		PayloadType:   MSG_SENT_SUCCESS,
		PayloadLength: MSG_SENT_SUCCESS_SIZE,
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

	testMessageSentSuccessPayload(result.Payload, expected.Payload, t)
}
