"use client";

import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import React, { useContext, useEffect, useRef, useState } from "react";
import {
  MessageSend,
  Packet,
  PayloadType,
  encodeMessageSend,
  encodePacket,
} from "@/app/lib/protocol";
import { SocketContext } from "@/app/lib/contexts";


export default function ChatInput() {
  const websocket = useContext(SocketContext);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [messages, setMessages] = useState<MessageSend[]>([]);

  function handleKeyDown(event: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (event.key === "Enter") {
      event.preventDefault();
      handleNewMessage();
    }
  }

  function handleNewMessage() {
    if (!websocket) {
      console.log("not connected");
      return
    }
    if (!textareaRef.current) {
      return;
    }

    const message = textareaRef.current.value;
    if (!message || message === "") {
      return;
    }

    const messageSend: MessageSend = {
      senderId: BigInt(15),
      recipientId: BigInt(16),
      sentAt: new Date(),
      content: message,
    };

    setMessages([...messages, messageSend]);

    textareaRef.current.value = "";

    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MessageSend,
      payloadLength: encodedMessage.length,
      payload: encodedMessage
    }
    const encodedPacket = encodePacket(packet);
    websocket.send(encodedPacket.buffer);
  }

  return (
    <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-zinc-800">
    <div className="flex-grow ">
        <textarea
        onKeyDown={handleKeyDown}
        ref={textareaRef}
        className="w-full p-2 bg-inherit text-white outline-none resize-none"
        placeholder="Type a message..."
        ></textarea>
    </div>
    <button onClick={handleNewMessage}>
        <div className="flex items-center px-2 py-2 mx-4 flex-shrink-0 rounded-xl bg-blue-500 hover:bg-blue-600 hover:cursor-pointer">
        <PaperAirplaneIcon className="w-6 h-6" />
        </div>
    </button>
    </div>
  );
}