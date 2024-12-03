
export interface Contact {
  contactId: bigint;
  name: string;
  email: string;
}

export interface User {
  id: BigInteger;
  firstName: string;
  lastName: string;
  email: string;
}

export interface ChatPreview {
  userId: bigint;
  userName: string;
  message: string;
  sentAt: Date;
}

export interface Message {
  id?: bigint;
  senderId: bigint;
  recipientId: bigint;
  sentAt?: Date;
  timestamp?: Date;
  content: string;
}

export interface Chat {
  contact: Contact;
  messages: Message[];
}