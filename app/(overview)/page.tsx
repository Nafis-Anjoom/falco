"use client";

import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox";
import { useEffect, useMemo, useRef, useState } from "react";
import { Contact, Message } from "../lib/definitions";
import ChatPane from "../ui/chat/chatpane";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/16/solid";
import {
  decodeMessageReceive,
  decodePacket,
  encodeMessageSend,
  encodePacket,
  Packet,
  PayloadType,
} from "../lib/protocol";
import Cookies from "js-cookie";

export default function Home() {
  const storedMessagesRef = useRef(new Map<number, Message[]>());
  const userIdRef = useRef(Number(Cookies.get("userId") ?? "0"));
  const [currentContact, setCurrentContact] = useState<Contact | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  console.log("messages: ", messages);

  useEffect(() => {
    if (!currentContact) {
      return;
    }

    console.log("current contact: ", currentContact);

    const messages = storedMessagesRef.current.get(currentContact.contactId);
    if (!messages) {
      //fetch the data
      console.log("fetching chat contactId: ", currentContact.contactId);
      let fetchedMessages: Message[] = [];

      storedMessagesRef.current.set(currentContact.contactId, fetchedMessages);
      setMessages(fetchedMessages);
    } else {
      setMessages(messages);
    }
  }, [currentContact]);

  const sendMessage = (content: string): void => {
    if (!currentContact) {
      return;
    }

    console.log("prepping message: ", currentContact.contactId, content);
    const messageSend: Message = {
      // the server will correct the senderId
      senderId: userIdRef.current,
      recipientId: currentContact.contactId,
      sentAt: new Date(),
      content: content,
    };

    const newMessages = [...messages, messageSend];
    setMessages(newMessages);
    storedMessagesRef.current.set(currentContact.contactId, newMessages);

    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MSG_SEND,
      payloadLength: encodedMessage.length,
      payload: encodedMessage,
    };
    const encodedPacket = encodePacket(packet);

    websocket.send(encodedPacket.buffer);
  };

  const websocket = useMemo(() => {
    const ws = new WebSocket("ws://localhost:3000/ws2");
    ws.onopen = () => {
      console.log("connected to message server");
    };

    ws.onerror = () => {
      ws.close();
      console.log("socket error. Closing connection.");
    };

    ws.onclose = () => {
      console.log("Disconnected from WebSocket");
    };

    return ws;
  }, []);

  websocket.onmessage = async (event) => {
    const blob: Blob = event.data;
    const arrayBuffer = await blob.arrayBuffer();
    const bytes = new Uint8Array(arrayBuffer);
    const packet = decodePacket(bytes);
    if (packet.payloadType == PayloadType.MSG_RECEIVE) {
      const messageReceive = decodeMessageReceive(packet.payload);

      // TODO: implement inbox system
      // When a message is received, adjust the inbox to display the chat on top
      // if the message thread is already in the storage, append the messages
      // else, just present the preview in the inbox
      console.log("messages on receipt: ", messages);
      console.log("current contact in focus: ", currentContact);
      if (!storedMessagesRef.current.has(messageReceive.senderId)) {
        console.log("msg received: ", messageReceive);
      } else {
        const newMessages = [...messages, messageReceive];
        storedMessagesRef.current.set(messageReceive.senderId, newMessages);
        if (
          currentContact &&
          currentContact.contactId === messageReceive.senderId
        ) {
          setMessages(newMessages);
        }
      }
    } else {
      console.log("payload type not supported now: ", packet);
    }
  };

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex h-full min-w-96">
        <ChatInbox setCurrentChat={setCurrentContact} />
      </div>
      {currentContact ? (
        <ChatPane
          userId={userIdRef.current}
          contact={currentContact}
          messages={messages}
          sendMessage={sendMessage}
        />
      ) : (
        <ChatPaneSkeleton />
      )}
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