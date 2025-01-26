"use client";

import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox-v2";
import { useEffect, useMemo, useRef, useState } from "react";
import { Contact, Message, MessageType, Packet, PayloadType } from "../lib/definitions_v2";
import ChatPane from "../ui/chat/chatpane";
import { ChatBubbleLeftRightIcon } from "@heroicons/react/16/solid";
import {
  decodeMessageReceive,
  decodeMessageSentSuccess,
  decodePacket,
  encodeMessageSend,
  encodePacket,
} from "../lib/protocol_v2";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import clsx from "clsx";

export default function Home() {
  const router = useRouter();
  const storedMessagesRef = useRef(new Map<number, Message[]>());
  const userIdRef = useRef(Number(Cookies.get("userId") ?? "0"));
  const [currentContact, setCurrentContact] = useState<Contact | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [isMessagesLoading, setIsMessagesLoading] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    (async () => {
      const response = await fetch("http://localhost:3000/validate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        credentials: "include",
      });

      if (!response.ok) {
        router.push("http://localhost:3001/login");
      } else {
        setIsLoggedIn(true);
      }
    })();
  }, []);

  useEffect(() => {
    if (!currentContact || !isLoggedIn) {
      return;
    }

    setIsMessagesLoading(true);

    const messages = storedMessagesRef.current.get(currentContact.contactId);
    if (!messages) {
      (async () => {
        try {
          const response = await fetch(
            `http://localhost:3000/chat/${currentContact.contactId}`,
            {
              method: "GET",
              headers: { "Content-Type": "application/json" },
              credentials: "include",
            }
          );

          if (response.ok) {
            const fetchedMessages: Message[] = await response.json();
            fetchedMessages.forEach((message) => {
              if (message.senderId === userIdRef.current) {
                message.type = MessageType.Send;
              } else {
                message.type = MessageType.Receive;
              }
              message.localUUID = crypto.randomUUID();
            });

            storedMessagesRef.current.set(
              currentContact.contactId,
              fetchedMessages
            );
            setMessages(fetchedMessages);
          }
        } catch (error) {
          console.log("error: ", error);
        }
      })();
    } else {
      setMessages(messages);
    }

    setIsMessagesLoading(false);
  }, [currentContact]);

  const sendMessage = (content: string): void => {
    if (!currentContact) {
      return;
    }

    const messageSend: Message = {
      // the server will correct the senderId
      type: MessageType.Send,
      senderId: userIdRef.current,
      recipientId: currentContact.contactId,
      sentAt: new Date(),
      content: content,
      localUUID: crypto.randomUUID()
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

  let websocket = useMemo(() => {
    let ws = new WebSocket("ws://localhost:3000/ws");
    ws.onopen = () => {
      setIsConnected(true);
      console.log("connected to message server");
    };

    ws.onerror = () => {
      ws.close(1000);
      setIsConnected(false);
      console.log("socket error. Attempting to reconnect.");
      setTimeout(() => {
      websocket = new WebSocket("ws://localhost:3000/ws");
        setTimeout(() => {
        if (websocket.readyState === WebSocket.OPEN) {
          console.log("reconnected");
          setIsConnected(true);
        } else {
          console.log("reconnected failed");
        }
        }, 2000);
      }, 3000);
    };

    ws.onclose = () => {
      setIsConnected(false);
      console.log("Disconnected from WebSocket. Attempting to reconnect...");
      setTimeout(() => {
      websocket = new WebSocket("ws://localhost:3000/ws");
        setTimeout(() => {
        if (websocket.readyState === WebSocket.OPEN) {
          console.log("reconnected");
          setIsConnected(true);
        } else {
          console.log("reconnection failed");
        }
        }, 2000);
      }, 3000);
    };

    return ws;
  }, []);

  function handleMessageReceive(payload: Uint8Array) {
    const messageReceive = decodeMessageReceive(payload);
    messageReceive.localUUID = crypto.randomUUID();

    // TODO: implement inbox system
    // When a message is received, adjust the inbox to display the chat on top
    // if the message thread is already in the storage, append the messages
    // else, just present the preview in the inbox
    if (storedMessagesRef.current.has(messageReceive.senderId)) {
      const currentMessages = storedMessagesRef.current.get(messageReceive.senderId) ?? [];
      const newMessages = [...currentMessages, messageReceive];
      storedMessagesRef.current.set(messageReceive.senderId, newMessages);
      if (
        currentContact && currentContact.contactId === messageReceive.senderId
      ) {
        setMessages(newMessages);
      }
    }
  }

    websocket.onmessage = async (event) => {
      const blob: Blob = event.data;
      const arrayBuffer = await blob.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const packet = decodePacket(bytes);
      switch (packet.payloadType) {
        case PayloadType.MSG_RECEIVE:
          handleMessageReceive(packet.payload);
          break;
        case PayloadType.MSG_SENT_SUCCESS:
          console.log("success: ", decodeMessageSentSuccess(packet.payload));
          break;
        default:
          console.log("not supported");
          break;
      }
    };
    
  return (
    <>
      {!isLoggedIn ? (
        <PageLoading />
      ) : (
        <div className="flex h-screen">
          <Sidebar />
          <div className="flex h-full min-w-96">
            <ChatInbox
              currentChat={currentContact}
              setCurrentChat={setCurrentContact}
            />
          </div>
          {currentContact ? (
            !isMessagesLoading ? (
              <ChatPane
                contact={currentContact}
                messages={messages}
                sendMessage={sendMessage}
              />
            ) : (
              <ChatPaneLoading />
            )
          ) : (
            <ChatPaneSkeleton />
          )}
          {!isConnected && 
            <div 
              className={clsx("absolute bg-red-600 px-4 py-1 rounded-xl bottom-3 left-14")}
            >
              Disconnected. Attempting to reconnect...
            </div>
          }
        </div>
      )}
    </>
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

function ChatPaneLoading() {
  return (
    <div className="flex h-screen">
      <div className="flex flex-col justify-center items-center w-full h-full">
        <ChatBubbleLeftRightIcon className="w-24 h-24" />
        <span className="font-semibold text-lg">Loading...</span>
      </div>
    </div>
  );
}

function PageLoading() {
  return (
    <div className="flex h-screen">
      <div className="flex h-full w-full items-center justify-center">
        <span className="text-xl">loading...</span>
      </div>
    </div>
  );
}