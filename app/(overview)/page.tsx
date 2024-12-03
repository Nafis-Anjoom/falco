"use client";

import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox";
import { useCallback, useMemo, useRef, useState } from "react";
import { Contact, Message } from "../lib/definitions";
import ChatPane from "../ui/chat/chatpane";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/16/solid";
import { decodeMessageReceive, decodePacket, encodeMessageSend, encodePacket, Packet, PayloadType } from "../lib/protocol";

const dummy1: Message[] = [
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

const dummy2: Message[] = [
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

interface ChatClient {
  currentUserId: bigint;
  websocket: WebSocket;
  isConnected: boolean;
  currentContact: Contact | null;
  messages: Message[];
  sendMessage: (content: string) => void;
  setCurrentChat: (contact: Contact | null) => void;
}

function useChatClient(): ChatClient {
  const [isConnected, setIsconnected] = useState(false);
  const storedMessagesRef = useRef(new Map<bigint, Message[]>());
  const [currentContact, setCurrentContact] = useState<Contact | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const currentUserId = BigInt(12);
  
  const websocket = useMemo(() => {
    const ws = new WebSocket("ws://localhost:3000/ws2");
    ws.onopen = () => {
      setIsconnected(true);
      console.log("connected to message server");
    };

    ws.onerror = () => {
      ws.close();
      console.log("socket error. Closing connection.");
    }

    ws.onclose = () => {
      setIsconnected(false);
      console.log("Disconnected from WebSocket");
    }

    ws.onmessage = async (event) => {
      const blob: Blob = event.data;
      const arrayBuffer = await blob.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const packet = decodePacket(bytes);
      const message = decodeMessageReceive(packet.payload);
      console.log("payload: ", packet.payload);
      console.log("message: ", message);
      return arrayBuffer;
    }

    return ws;
  }, []);

  const sendMessage = (content: string): void => {
    if (!currentContact) {
      return;
    }

    console.log("prepping message: ", currentContact.contactId, content);
    const messageSend: Message = {
      senderId: BigInt(15),
      recipientId: BigInt(currentContact.contactId),
      sentAt: new Date(),
      content: content,
    };

    const newMessages = [...messages, messageSend];
    setMessages(newMessages);
    storedMessagesRef.current.set(currentContact.contactId, newMessages);

    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MessageSend,
      payloadLength: encodedMessage.length,
      payload: encodedMessage,
    };
    const encodedPacket = encodePacket(packet);

    websocket.send(encodedPacket.buffer);
  }

  const setCurrentChat = useCallback((contact: Contact | null): void => {
    // if contact in stored chat, then change state
    // elese, fetch from server, store in map, then change state
    if (!contact) {
      return;
    }

    const messages = storedMessagesRef.current.get(contact.contactId);
    if (!messages) {
      //fetch the data
      console.log("fetching chat contactId: ", contact.contactId);
      let fetchedMessages: Message[] = [];
      if (BigInt(contact.contactId) === BigInt(18)) {
        fetchedMessages = dummy2;
      } else {
        fetchedMessages = dummy1;
      }

      storedMessagesRef.current.set(contact.contactId, fetchedMessages);
      setCurrentContact(contact);
      setMessages(fetchedMessages);
    } else {
      setCurrentContact(contact);
      setMessages(messages);
    }
  }, []);

  console.log("use chat client");

  return {
    currentUserId: currentUserId,
    websocket: websocket,
    isConnected: isConnected,
    currentContact: currentContact,
    messages: messages,
    sendMessage: sendMessage,
    setCurrentChat: setCurrentChat,
  }
}

export default function Home() {
  const chatClient = useChatClient();

  return (
      <div className="flex h-screen">
        <Sidebar />
        <div className="flex h-full min-w-96">
          <ChatInbox setCurrentChat={chatClient.setCurrentChat} />
        </div>
        {chatClient.currentContact ? 
          <ChatPane 
            contact={chatClient.currentContact}
            messages={chatClient.messages} 
            sendMessage={chatClient.sendMessage}
          /> 
          : <ChatPaneSkeleton />}
      </div>
  );
}

function ChatPaneSkeleton() {
    return (
      <div className="flex flex-col justify-center items-center w-full h-full">
        <ChatBubbleLeftRightIcon className="w-24 h-24" />
        <span className="font-semibold text-lg">Start a chat</span>
      </div>
    );
}