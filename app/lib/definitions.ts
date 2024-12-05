export interface Contact {
  contactId: number;
  name: string;
  email: string;
}

export interface User {
  id: number;
  firstName: string;
  lastName: string;
  email: string;
}

export interface ChatPreview {
  userId: number;
  userName: string;
  message: string;
  sentAt: Date;
}

export interface Message {
  id?: bigint;
  senderId: number;
  recipientId: number;
  sentAt?: Date;
  timestamp?: Date;
  content: string;
}

export interface Chat {
  contact: Contact;
  messages: Message[];
}

export interface Packet {
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

export interface MessageSentSuccess {
  messageId: bigint,
  recipientId: number,
  timestamp: Date,
  sentAt: Date
}

