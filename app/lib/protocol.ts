import { Message } from "./definitions";

export enum PayloadType {
	MSG_READ_SUCCESS,
	MSG_READ_FAIL,
	MSG_DELV_SUCCESS,
	MSG_DELV_FAIL,
	MSG_SENT_SUCCESS,
	MSG_SENT_FAIL,
    MSG_SEND,
    MSG_RECEIVE,
    SYNC_THREAD,
	CONN_ERR,
	CONN_FIN,
	CONN_INIT
}

export interface Packet {
    version: number; // uint8
    payloadType: PayloadType; //uint8
    payloadLength: number; // uint16
    payload: Uint8Array; // Byte array
}

export function encodeMessageSend(message: Message): Uint8Array {
    const encoder = new TextEncoder();
    const contentBytes = encoder.encode(message.content);

    const buffer = new ArrayBuffer(24 + contentBytes.length);
    const view = new DataView(buffer);

    view.setBigUint64(0, BigInt(message.senderId), false);
    view.setBigUint64(8, BigInt(message.recipientId), false);

    // timestamp could be null
    const timestamp = message.sentAt?.getTime();
    view.setBigUint64(16, BigInt(timestamp ?? 0), false);

    const contentOffset = 24;
    const output = new Uint8Array(buffer);
    output.set(contentBytes, contentOffset);

    return output;
}

export function encodeMessageSendPacket(message: Message): Uint8Array {
    const encoder = new TextEncoder();
    const contentBytes = encoder.encode(message.content);

    const buffer = new ArrayBuffer(24 + contentBytes.length);
    const view = new DataView(buffer);

    view.setBigUint64(0, BigInt(message.senderId), false);
    view.setBigUint64(8, BigInt(message.recipientId), false);
    // TODO: investigate default values
    const timestamp = message.sentAt?.getTime();
    view.setBigUint64(16, BigInt(timestamp ?? 0), false);

    const contentOffset = 24;
    const payload = new Uint8Array(buffer);
    payload.set(contentBytes, contentOffset);

    const packet: Packet = {
        version: 1,
        payloadType: PayloadType.MSG_SEND,
        payloadLength: payload.length,
        payload: payload
    }

    return encodePacket(packet);
}

export function decodeMessageSend(buffer: Uint8Array): Message {
    const view = new DataView(buffer.buffer);

    const senderId = Number(view.getBigUint64(0, true));
    const recipientId = Number(view.getBigUint64(8, true));

    const timestamp = view.getUint32(16, true);
    const sentAt = new Date(timestamp);

    const content = new TextDecoder().decode(buffer.subarray(20));

    return { senderId, recipientId, sentAt, content };
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
        id,
        senderId,
        recipientId,
        timestamp,
        content
    }
}

export function decodeMessageReceivePacket(bytes: Uint8Array): Message {
    const view = new DataView(bytes.buffer);
    
    const protocolVersion = view.getUint8(0);
    const payloadType = view.getUint8(1);
    const payloadLength = view.getUint16(2, false);
    
    const id = view.getBigInt64(4, false);
    const senderId = Number(view.getBigInt64(12, false));
    const recipientId = Number(view.getBigInt64(20, false));
    const timestamp = new Date(Number(view.getBigInt64(28, false)));
    const content = new TextDecoder().decode(bytes.subarray(36));

    return {
        id,
        senderId,
        recipientId,
        timestamp,
        content
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