package protocol_v2

import (
	"testing"
	"time"
)

func TestMessageReceiveEncoding(t *testing.T) {
	messageIn := MessageReceieve{
		Id:          1,
		SenderId:    2,
		RecipientId: 3,
		Timestamp:   time.Unix(12345678, 0),
		Content:     "hello",
	}

	result, _ := messageIn.MarshalBinary()
	expected := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
		0x68, 0x65, 0x6C, 0x6C, 0x6F,
	}

	testMessageReceivePayload(result, expected, t)
}

func TestMessageReceiveDecoding(t *testing.T) {
	input := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
		0x68, 0x65, 0x6C, 0x6C, 0x6F,
	}

	var result MessageReceieve
	err := result.UnmarshalBinary(input)
	if err != nil {
		t.Fatalf("unexpected error during unmarshal: %v", err)
	}

	expected := MessageReceieve{
		Id:          1,
		SenderId:    2,
		RecipientId: 3,
		Timestamp:   time.Unix(12345678, 0),
		Content:     "hello",
	}

	if result != expected {
		t.Fatalf("Expected: %+v. Result: %+v", expected, result)
	}
}

func TestMessageSendEncoding(t *testing.T) {
	messageIn := MessageSend{
		SenderId:    2,
		RecipientId: 3,
		SentAt:      time.Unix(12345678, 0),
		LocalUUID:   "f507e7f6-8b0c-4fde-9e62-f1f0260b9a38",
		Content:     "hello",
	}

	result, _ := messageIn.MarshalBinary()
	expected := []byte{
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0xBC, 0x61, 0x4E,
        0x66, 0x35, 0x30, 0x37, 0x65, 0x37, 0x66, 0x36, 0x2d, 0x38, 0x62, 0x30, 0x63, 0x2d, 0x34, 0x66, 0x64, 0x65, 0x2d, 0x39, 0x65, 0x36, 0x32, 0x2d, 0x66, 0x31, 0x66, 0x30, 0x32, 0x36, 0x30, 0x62, 0x39, 0x61, 0x33, 0x38,
		0x68, 0x65, 0x6C, 0x6C, 0x6F,
	}

	testMessageSendPayload(expected, result, t)
}

// utils
func testMessageSendPayload(expected []byte, result []byte, t *testing.T) {
	if len(expected) != len(result) {
		t.Fatalf("Expected size %d. Got %d", len(expected), len(result))
	}

    for i := 0; i < 8; i++ {
		if expected[i] != result[i] {
			t.Errorf("Incorrect SenderId encoding. Result: % x. Expected % x\n", expected[0:8], result[0:8])
            break
		}
	}

    for i := 8; i < 16; i++ {
		if expected[i] != result[i] {
            t.Errorf("Incorrect RecipientId encoding. Result: % x. Expected % x\n", expected[8:16], result[8:16])
            break
		}
	}

    for i := 16; i < 24; i++ {
		if expected[i] != result[i] {
            t.Errorf("Incorrect SentAt encoding. Result: % x. Expected % x\n", expected[16:24], result[16:24])
            break
		}
	}

    for i := 24; i < 60; i++ {
		if expected[i] != result[i] {
            t.Errorf("Incorrect localUUID encoding. Result: % x. Expected % x\n", expected[24:60], result[24:60])
            break
		}
	}

    for i := 60; i < len(expected); i++ {
		if expected[i] != result[i] {
            t.Errorf("Incorrect Content encoding. Result: % x. Expected % x\n", expected[60:], result[60:])
            break
		}
	}
}

func testMessageReceivePayload(b1 []byte, b2 []byte, t *testing.T) {
	if len(b1) != len(b2) {
		t.Fatalf("Expected size 37. Got %d", len(b1))
	}

    for i := 0; i < 8; i++ {
		if b1[i] != b2[i] {
			t.Errorf("Incorrect Id encoding. Result: % x. Expected % x\n", b1[0:8], b2[0:8])
            break
		}
	}

    for i := 8; i < 16; i++ {
		if b1[i] != b2[i] {
            t.Errorf("Incorrect SenderId encoding. Result: % x. Expected % x\n", b1[8:16], b2[8:16])
            break
		}
	}

    for i := 16; i < 24; i++ {
		if b1[i] != b2[i] {
            t.Errorf("Incorrect Recipient encoding. Result: % x. Expected % x\n", b1[16:24], b2[16:24])
            break
		}
	}

    for i := 24; i < 32; i++ {
		if b1[i] != b2[i] {
            t.Errorf("Incorrect Timestamp encoding. Result: % x. Expected % x\n", b1[24:32], b2[24:32])
            break
		}
	}

    for i := 32; i < len(b1); i++ {
		if b1[i] != b2[i] {
            t.Errorf("Incorrect Content encoding. Result: % x. Expected % x\n", b1[32:37], b2[32:37])
            break
		}
	}
}

func testMessageSentSuccessPayload(b1 []byte, b2 []byte, t *testing.T) {
	if len(b1) != len(b2) {
		t.Fatalf("Expected size 37. Got %d", len(b1))
	}

    for i := 0; i < 8; i++ {
		if b1[i] != b2[i] {
			t.Errorf("Incorrect MessageId encoding. Result: % x. Expected % x\n", b1[0:8], b2[0:8])
            break
		}
	}

    for i := 8; i < 16; i++ {
		if b1[i] != b2[i] {
            t.Errorf("Incorrect Timestamp encoding. Result: % x. Expected % x\n", b1[8:16], b2[8:16])
            break
		}
	}
}
