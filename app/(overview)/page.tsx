"use client";

import Sidebar from "@/app/ui/sidebar";
import ChatInbox from "@/app/ui/chat/chatInbox";
import ChatPane from "../ui/chat/chatpane";

export default function Home() {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <ChatInbox />
      <ChatPane />
    </div>
  );
}