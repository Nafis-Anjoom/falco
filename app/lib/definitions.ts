
export interface Contact {
  contactId: number;
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
  userId: BigInt;
  userName: string;
  message: string;
  sentAt: Date;
}