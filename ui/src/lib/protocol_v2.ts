import { Message, MessageSentSuccess, MessageType, Packet } from "./definitions_v2";

const MESSAGE_SEND_HEADER_SIZE = 60;

export function encodeMessageSend(message: Message): Uint8Array {
    const encoder = new TextEncoder();

    // const buffer = new ArrayBuffer(MESSAGE_SEND_HEADER_SIZE + contentBytes.length);
    const buffer = new ArrayBuffer(MESSAGE_SEND_HEADER_SIZE + message.content.length);
    const view = new DataView(buffer);

    view.setBigUint64(0, BigInt(message.senderId), false);
    view.setBigUint64(8, BigInt(message.recipientId), false);

    if (!message.sentAt) {
        message.sentAt = new Date();
    }
    const timestamp = Math.floor(message.sentAt.getTime() / 1000);
    view.setBigUint64(16, BigInt(timestamp), false);

    const output = new Uint8Array(buffer);

    const localUUIDOffset = 24;
    const localUUIDBytes = encoder.encode(message.localUUID);
    output.set(localUUIDBytes, localUUIDOffset);

    const contentOffset = 60;
    const contentBytes = encoder.encode(message.content);
    output.set(contentBytes, contentOffset);

    return output;
}

export function decodeMessageSend(buffer: Uint8Array): Message {
    const view = new DataView(buffer.buffer);

    const senderId = Number(view.getBigUint64(0, true));
    const recipientId = Number(view.getBigUint64(8, true));

    const timestamp = view.getUint32(16, true);
    const sentAt = new Date(timestamp);

    const content = new TextDecoder().decode(buffer.subarray(20));

    return {
        type: MessageType.Send,
        senderId,
        recipientId,
        sentAt,
        content
    };
}

export function decodeMessageReceive(bytes: Uint8Array): Message {
    // the underlying array buffer contains the whole packet
    // the offset of 4 ignores the packet header
    const view = new DataView(bytes.buffer, 4);

    const id = view.getBigInt64(0, false);
    const senderId = Number(view.getBigInt64(8, false));
    const recipientId = Number(view.getBigInt64(16, false));
    const timestamp = new Date(Number(view.getBigInt64(24, false)) * 1000);
    const content = new TextDecoder().decode(bytes.subarray(32));

    return {
        type: MessageType.Receive,
        id,
        senderId,
        recipientId,
        timestamp,
        content
    }
}

export function decodeMessageSentSuccess(bytes: Uint8Array): MessageSentSuccess {
    const view = new DataView(bytes.buffer, 4);

    const messageId = view.getBigInt64(0, false);
    const recipientId = Number(view.getBigInt64(8, false));
    const timestamp = new Date(Number(view.getBigInt64(16, false)) * 1000);
    const localUUID = new TextDecoder().decode(bytes.subarray(24));

    return {
        messageId: messageId,
        recipientId: recipientId,
        timestamp: timestamp,
        localUUID: localUUID
    }
}

export function encodePacket(packet: Packet): Uint8Array {
    const buffer = new ArrayBuffer(4 + packet.payload.byteLength);
    const view = new DataView(buffer);

    view.setUint8(0, packet.version);
    view.setUint8(1, packet.payloadType);
    view.setUint16(2, packet.payloadLength, false);

    const uint8Array = new Uint8Array(buffer);
    uint8Array.set(packet.payload, 4);

    return uint8Array;
}

export function decodePacket(buffer: Uint8Array): Packet {
    const view = new DataView(buffer.buffer);

    const version = view.getUint8(0);
    const payloadType = view.getUint8(1);
    const payloadLength = view.getUint16(2, false);

    const payload = buffer.subarray(4);

    return { version, payloadType, payloadLength, payload };
}
