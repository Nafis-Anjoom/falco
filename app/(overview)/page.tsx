"use client";

import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox";
import { useCallback, useMemo, useRef, useState } from "react";
import { Contact, Message } from "../lib/definitions";
import ChatPane from "../ui/chat/chatpane";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/16/solid";
import {
  decodeMessageReceive,
  decodePacket,
  encodeMessageSend,
  encodePacket,
  Packet,
  PayloadType
} from "../lib/protocol";

interface ChatClient {
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
      const messageReceive = decodeMessageReceive(packet.payload);
      
      // TODO: implement inbox system
      // When a message is received, adjust the inbox to display the chat on top
      // if the message thread is already in the storage, append the messages
      // else, just present the preview in the inbox
      if (!storedMessagesRef.current.has(messageReceive.senderId)) {
        console.log("msg received: ", messageReceive);
      } else {
        const newMessages = [...messages, messageReceive];
        storedMessagesRef.current.set(messageReceive.senderId, newMessages);
        if (currentContact && BigInt(currentContact.contactId) === BigInt(messageReceive.senderId)) {
          setMessages(newMessages);
        }
      }

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
      // the server will correct the senderId
      senderId: BigInt(0),
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