
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