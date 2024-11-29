"use client";

import ChatInbox from "../ui/chat/chatInbox";
import Sidebar from "../ui/sidebar";
import { SocketContext } from "@/app/lib/contexts"

const websocket = new WebSocket("ws://localhost:3000/ws2");
websocket.onopen = () => {
  console.log("connected to message server");
};

websocket.onerror = () => {
  console.log("socket error");
}

websocket.onclose = () => {
  console.log("Disconnected from WebSocket");
}

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <SocketContext.Provider value={websocket}>
      <div className="flex h-screen">
        <Sidebar />
        <div className="flex h-full min-w-96">
          <ChatInbox />
        </div>
        {children}
      </div>
    </SocketContext.Provider>
  );
}