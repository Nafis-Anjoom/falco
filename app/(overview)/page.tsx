"use client";

import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox";
import { useCallback, useMemo, useRef, useState } from "react";
import { Chat, Contact, Message } from "../lib/definitions";
import ChatPane from "../ui/chat/chatpane";

type Socket = {
  websocket: WebSocket;
  isConnected: boolean;
}

function useWebsocket(): Socket {
  const [isConnected, setIsconnected] = useState(false);
  
  const websocket = new WebSocket("ws://localhost:3000/ws2");
  websocket.onopen = () => {
    setIsconnected(true);
    console.log("connected to message server");
  };

  websocket.onerror = () => {
    websocket.close();
    console.log("socket error. Closing connection.");
  }

  websocket.onclose = () => {
    setIsconnected(false);
    console.log("Disconnected from WebSocket");
  }
  
  return { websocket, isConnected }
}

interface ChatClient {
  currentUserId: bigint;
  websocket: WebSocket;
  isConnected: boolean;
  currentChat: Chat | null;
  setCurrentChat: (contact: Contact | null) => void;
}

function useChatClient(): ChatClient {
  const [isConnected, setIsconnected] = useState(false);

  // const storedChats = new Map<bigint, Chat>();
  const storedChatsRef = useRef(new Map<bigint, Chat>());
  const [currentChat, changeChat] = useState<Chat | null>(null);
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

    return ws;
  }, []);

  const setCurrentChat = useCallback((contact: Contact | null): void => {
    // if contact in stored chat, then change state
    // elese, fetch from server, store in map, then change state
    if (!contact) {
      return;
    }

    const chat = storedChatsRef.current.get(contact.contactId);
    if (!chat) {
      //fetch the data
      console.log("fetching chat");
      const fetchedChat: Chat = {
        contact: contact,
        messages: []
      }

      storedChatsRef.current.set(contact.contactId, fetchedChat);
      changeChat(fetchedChat);
    } else {
      changeChat(chat);
    }
  }, []);

  console.log("use chat client");

  return {
    currentUserId: currentUserId,
    websocket: websocket,
    isConnected: isConnected,
    currentChat: currentChat,
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
        <ChatPane chat={chatClient.currentChat} />
      </div>
  );
}