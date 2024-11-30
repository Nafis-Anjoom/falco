"use client";

import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/solid";
import Sidebar from "../ui/sidebar";
import ChatInbox from "../ui/chat/chatInbox";
import { useState } from "react";
import { Contact } from "../lib/definitions";
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

export default function Home() {
  const socket = useWebsocket();
  const [contact, setContact] = useState<null | Contact>(null);

  return (
      <div className="flex h-screen">
        <Sidebar />
        <div className="flex h-full min-w-96">
          <ChatInbox setChat={setContact} />
        </div>
        {!contact ? <ChatPaneHome /> : <ChatPane contact={contact} />}
      </div>
  );
}

function ChatPaneHome() {
  return (
    <div className="flex flex-col justify-center items-center w-full h-full">
      <ChatBubbleLeftRightIcon className="w-24 h-24" />
      <span className="font-semibold text-lg">Start a chat</span>
    </div>
  );
}