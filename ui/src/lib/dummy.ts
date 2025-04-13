import { Message } from "./definitions_v1";

export const dummy1: Message[] = [
  {
    senderId: BigInt(15),
    recipientId: BigInt(17),
    content: "Hey! How are you doing?",
    sentAt: new Date('2024-11-30T10:15:00Z'),
  },
  {
    senderId: BigInt(17),
    recipientId: BigInt(15),
    content: "I'm good, thanks! How about you?",
    sentAt: new Date('2024-11-30T10:16:00Z'),
  },
  {
    senderId: BigInt(15),
    recipientId: BigInt(17),
    content: "Doing well, thanks for asking!",
    sentAt: new Date('2024-11-30T10:17:00Z'),
  },
  {
    senderId: BigInt(17),
    recipientId: BigInt(15),
    content: "Great to hear! Have any plans for today?",
    sentAt: new Date('2024-11-30T10:18:00Z'),
  },
  {
    senderId: BigInt(15),
    recipientId: BigInt(17),
    content: "Not much, just planning to relax. You?",
    sentAt: new Date('2024-11-30T10:19:00Z'),
  },
  {
    senderId: BigInt(17),
    recipientId: BigInt(15),
    content: "Same here, maybe catch up on some shows.",
    sentAt: new Date('2024-11-30T10:20:00Z'),
  },
];

export const dummy2: Message[] = [
  {
    senderId: BigInt(15),
    recipientId: BigInt(18),
    content: "Hi there! Are you free to chat?",
    sentAt: new Date('2024-11-30T14:00:00Z'),
  },
  {
    senderId: BigInt(18),
    recipientId: BigInt(15),
    content: "Hey! Sure, what's up?",
    sentAt: new Date('2024-11-30T14:01:30Z'),
  },
  {
    senderId: BigInt(15),
    recipientId: BigInt(18),
    content: "I wanted to discuss the project deadline. Can we extend it?",
    sentAt: new Date('2024-11-30T14:02:45Z'),
  },
  {
    senderId: BigInt(18),
    recipientId: BigInt(15),
    content: "That might be possible. How much more time do you need?",
    sentAt: new Date('2024-11-30T14:03:20Z'),
  },
  {
    senderId: BigInt(15),
    recipientId: BigInt(18),
    content: "Just a couple of extra days should be enough.",
    sentAt: new Date('2024-11-30T14:04:10Z'),
  },
  {
    senderId: BigInt(18),
    recipientId: BigInt(15),
    content: "Alright, let me check with the team. I'll get back to you soon.",
    sentAt: new Date('2024-11-30T14:05:00Z'),
  },
];