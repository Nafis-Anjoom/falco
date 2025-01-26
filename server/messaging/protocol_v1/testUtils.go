package protocol_v1

import "testing"

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
