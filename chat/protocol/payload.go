package protocol

import (
	"encoding"
	"encoding/binary"
	"time"
)

const (
	MSG_RECEIVE_HEADER_SIZE = 32
	MSG_SEND_HEADER_SIZE    = 24
	MSG_STATUS_SENT_SIZE    = 16
)

type Payload interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
	Type() PayloadType
	Length() int
}

type MessageReceieve struct {
	Id          int64     `json:"id"`
	SenderId    int64     `json:"sender_id"`
	RecipientId int64     `json:"recipient_id"`
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
	SenderId    int64     `json:"sender_id"`
	RecipientId int64     `json:"recipient_id"`
	SentAt      time.Time `json:"sent_at"`
	Content     string    `json:"content"`
}

func (ms *MessageSend) MarshalBinary() (data []byte, err error) {
	buffer := make([]byte, MSG_SEND_HEADER_SIZE+len(ms.Content))

	binary.BigEndian.PutUint16(buffer[0:], uint16(ms.SenderId))
	binary.BigEndian.PutUint16(buffer[8:], uint16(ms.RecipientId))
	binary.BigEndian.PutUint16(buffer[16:], uint16(ms.SentAt.Unix()))
	copy(buffer[MSG_SEND_HEADER_SIZE:], ms.Content)

	return buffer, nil
}

func (ms *MessageSend) UnmarshalBinary(data []byte) error {
	ms.SenderId = int64(binary.BigEndian.Uint64(data[0:]))
	ms.RecipientId = int64(binary.BigEndian.Uint64(data[8:]))
	sentAt := int64(binary.BigEndian.Uint64(data[16:]))
	ms.SentAt = time.Unix(sentAt, 0)
	ms.Content = string(data[MSG_SEND_HEADER_SIZE:])

	return nil
}

func (ms *MessageSend) Type() PayloadType {
	return MSG_SEND
}

func (ms *MessageSend) Length() int {
	return MSG_SEND_HEADER_SIZE + len(ms.Content)
}

type MessageStatusSent struct {
	MessageId int64     `json:"message_id"`
	Timestamp time.Time `json:"timestamp"`
}

func (ms *MessageStatusSent) MarshalBinary() (data []byte, err error) {
	buffer := make([]byte, MSG_STATUS_SENT_SIZE)

	binary.BigEndian.PutUint64(buffer[0:], uint64(ms.MessageId))
	binary.BigEndian.PutUint64(buffer[8:], uint64(ms.Timestamp.Unix()))

	return buffer, nil
}

func (ms *MessageStatusSent) UnmarshalBinary(data []byte) error {
	ms.MessageId = int64(binary.BigEndian.Uint64(data[0:]))
	timestamp := int64(binary.BigEndian.Uint64(data[8:]))
	ms.Timestamp = time.Unix(timestamp, 0)

	return nil
}

func (ms *MessageStatusSent) Type() PayloadType {
	return MSG_STATUS_SENT
}

func (ms *MessageStatusSent) Length() int {
	return MSG_STATUS_SENT_SIZE
}
