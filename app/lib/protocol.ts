export enum PayloadType {
    MessageSend = 6,
    MessageReceive = 2
}

export interface Packet {
    version: number; // uint8
    payloadType: PayloadType; //uint8
    payloadLength: number; // uint16
    payload: Uint8Array; // Byte array
}

export interface MessageSend {
    senderId: bigint; // int64
    recipientId: bigint; // int64
    sentAt: Date; // int64
    content: string; //Byte array
}

// Helper function to encode a MessageSend object into a buffer
export function encodeMessageSend(message: MessageSend): Uint8Array {
    const encoder = new TextEncoder();
    const contentBytes = encoder.encode(message.content);

    const buffer = new ArrayBuffer(24 + contentBytes.length);
    const view = new DataView(buffer);

    // Serialize senderId and recipientId (int64 -> BigInt64Array)
    view.setBigUint64(0, message.senderId, false);
    view.setBigUint64(8, message.recipientId, false);

    // Serialize sentAt (timestamp in milliseconds)
    const timestamp = message.sentAt.getTime();
    view.setBigUint64(16, BigInt(timestamp), false);

    // Serialize content
    const contentOffset = 24; // Adjust based on timestamp size
    const uint8Array = new Uint8Array(buffer);
    uint8Array.set(contentBytes, contentOffset);

    // return uint8Array;
    // const uint8Array = new Uint8Array(buffer);
    return uint8Array;
}

// Helper function to decode a buffer into a MessageSend object
export function decodeMessageSend(buffer: Uint8Array): MessageSend {
    const view = new DataView(buffer.buffer);

    const senderId = view.getBigUint64(0, true);
    const recipientId = view.getBigUint64(8, true);

    const timestamp = view.getUint32(16, true);
    const sentAt = new Date(timestamp);

    const content = new TextDecoder().decode(buffer.subarray(20));

    return { senderId, recipientId, sentAt, content };
}

// Function to encode a Packet
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

// Function to decode a buffer into a Packet
function decodePacket(buffer: Uint8Array): Packet {
    const view = new DataView(buffer.buffer);

    const version = view.getUint8(0);
    const payloadType = view.getUint8(1);
    const payloadLength = view.getUint16(2, false);

    const payload = buffer.subarray(4, 4 + payloadLength);

    return { version, payloadType, payloadLength, payload };
}