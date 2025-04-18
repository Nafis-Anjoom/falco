package protocol_v2

import (
	"encoding"
	"encoding/binary"
	"time"
)

const (
	MSG_RECEIVE_HEADER_SIZE = 32
	MSG_SEND_HEADER_SIZE    = 60
	MSG_SENT_SUCCESS_SIZE   = 60
	SYNC_THREAD_SIZE        = 16
)

type Payload interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Type() PayloadType
	Length() int
}

type MessageReceieve struct {
	Id          int64     `json:"id"`
	SenderId    int64     `json:"senderId"`
	RecipientId int64     `json:"recipientId"`
	Timestamp   time.Time `json:"timestamp"`
	Content     string    `json:"content"`
}

func (mr *MessageReceieve) MarshalBinary() (data []byte, err error) {
	buffer := make([]byte, 32+len(mr.Content))

	binary.BigEndian.PutUint64(buffer[0:], uint64(mr.Id))
	binary.BigEndian.PutUint64(buffer[8:], uint64(mr.SenderId))
	binary.BigEndian.PutUint64(buffer[16:], uint64(mr.RecipientId))
	binary.BigEndian.PutUint64(buffer[24:], uint64(mr.Timestamp.Unix()))
	copy(buffer[32:], mr.Content)

	return buffer, nil
}

func (mr *MessageReceieve) UnmarshalBinary(data []byte) error {
	mr.Id = int64(binary.BigEndian.Uint64(data[0:]))
	mr.SenderId = int64(binary.BigEndian.Uint64(data[8:]))
	mr.RecipientId = int64(binary.BigEndian.Uint64(data[16:]))
	timestamp := int64(binary.BigEndian.Uint64(data[24:]))
	mr.Timestamp = time.Unix(timestamp, 0)
	mr.Content = string(data[32:])

	return nil
}

func (mr *MessageReceieve) Type() PayloadType {
	return MSG_RECEIVE
}

func (mr *MessageReceieve) Length() int {
	return MSG_RECEIVE_HEADER_SIZE + len(mr.Content)
}

type MessageSend struct {
	SenderId    int64     `json:"senderId"`
	RecipientId int64     `json:"recipientId"`
	SentAt      time.Time `json:"sentAt"`
	LocalUUID   string    `json:"localUUID"` // 36 characters long uuid
	Content     string    `json:"content"`
}

func (ms *MessageSend) MarshalBinary() (data []byte, err error) {
	buffer := make([]byte, MSG_SEND_HEADER_SIZE+len(ms.Content))

	binary.BigEndian.PutUint64(buffer[0:], uint64(ms.SenderId))
	binary.BigEndian.PutUint64(buffer[8:], uint64(ms.RecipientId))
	binary.BigEndian.PutUint64(buffer[16:], uint64(ms.SentAt.Unix()))
	copy(buffer[24:60], ms.LocalUUID)
	copy(buffer[MSG_SEND_HEADER_SIZE:], ms.Content)

	return buffer, nil
}

func (ms *MessageSend) UnmarshalBinary(data []byte) error {
	ms.SenderId = int64(binary.BigEndian.Uint64(data[0:]))
	ms.RecipientId = int64(binary.BigEndian.Uint64(data[8:]))
	sentAt := int64(binary.BigEndian.Uint64(data[16:]))
	ms.SentAt = time.Unix(sentAt, 0)
	ms.LocalUUID = string(data[24:60])
	ms.Content = string(data[MSG_SEND_HEADER_SIZE:])

	return nil
}

func (ms *MessageSend) Type() PayloadType {
	return MSG_SEND
}

func (ms *MessageSend) Length() int {
	return MSG_SEND_HEADER_SIZE + len(ms.Content)
}

type MessageSentSuccess struct {
	MessageId   int64     `json:"messageId"`
	RecipientId int64     `json:"recipientId"`
	Timestamp   time.Time `json:"timestamp"`
	LocalUUID   string    `json:"localUUID"`
}

func (ms *MessageSentSuccess) MarshalBinary() (data []byte, err error) {
	buffer := make([]byte, MSG_SENT_SUCCESS_SIZE)

	binary.BigEndian.PutUint64(buffer[0:], uint64(ms.MessageId))
	binary.BigEndian.PutUint64(buffer[8:], uint64(ms.RecipientId))
	binary.BigEndian.PutUint64(buffer[16:], uint64(ms.Timestamp.Unix()))
	copy(buffer[24:], ms.LocalUUID)

	return buffer, nil
}

func (ms *MessageSentSuccess) UnmarshalBinary(data []byte) error {
	ms.MessageId = int64(binary.BigEndian.Uint64(data[0:]))
	ms.RecipientId = int64(binary.BigEndian.Uint64(data[8:]))
	timestamp := int64(binary.BigEndian.Uint64(data[16:]))
	ms.Timestamp = time.Unix(timestamp, 0)
	ms.LocalUUID = string(data[24:])

	return nil
}

func (ms *MessageSentSuccess) Type() PayloadType {
	return MSG_SENT_SUCCESS
}

func (ms *MessageSentSuccess) Length() int {
	return MSG_SENT_SUCCESS_SIZE
}
