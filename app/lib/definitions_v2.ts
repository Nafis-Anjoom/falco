export type Contact = {
    contactId: number;
    name: string;
    email: string;
}

export type User = {
    id: number;
    firstName: string;
    lastName: string;
    email: string;
}

export type ChatPreview = {
    userId: number;
    userName: string;
    message: string;
    sentAt: Date;
}

export type Message = {
    id?: bigint;
    senderId: number;
    recipientId: number;
    sentAt?: Date | undefined;
    deliveredAt?: Date | undefined;
    seenAt?: Date | undefined;
    timestamp?: Date | undefined;
    content: string;
}

export type MessageSend = {
    senderId: number;
    recipientId: number;
    localId: number;
    sentAt: Date;
    content: string;
}

export type MessageReceive = {
    id: bigint;
    senderId: number;
    recipientId: number;
    timestamp: Date;
    content: string;
}

export type MessageSentSuccess = {
    messageId: bigint;
    localId: string;
    timestamp: Date;
}

export type Chat = {
    contact: Contact;
    messages: Message[];
}

export type Packet = {
    version: number; // uint8
    payloadType: PayloadType; //uint8
    payloadLength: number; // uint16
    payload: Uint8Array; // Byte array
}

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