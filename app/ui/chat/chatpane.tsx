"use client";

import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import Message from "./message";
import React, { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { MessageSend, Packet, PayloadType, encodeMessageSend, encodePacket } from "@/app/lib/protocol";

const socketURL = "ws://localhost:3000/ws2";

export default function ChatPane() {
  const router = useRouter();
  const websocketRef = useRef<WebSocket| null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const token = localStorage.getItem("socketToken");
    if (!token || token == null) {
      router.push("/login");
    }

    websocketRef.current = new WebSocket(`${socketURL}?token=${token}`);
    websocketRef.current.onopen = () => {
      console.log("connected to message server");
      setIsConnected(true);
    };

    websocketRef.current.onclose = () => {
      setIsConnected(false);
      console.log('Disconnected from WebSocket');
    };
  }, []);

  function handleNewMessage() {
    const message = textareaRef.current?.value;
    if (!message) {
      return;
    }
    const messageSend: MessageSend = {
      senderId: BigInt(5),
      recipientId: BigInt(1),
      sentAt: new Date(),
      content: message
    }
    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MessageSend,
      payloadLength: encodedMessage.length,
      payload: encodedMessage
    }
    const encodedPacket = encodePacket(packet);
    console.log(encodedPacket);
    
    websocketRef.current?.send(encodedPacket.buffer);
  }

  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      <div className="flex flex-shrink-0 flex-grow-0 border-b-2 border-blue-500 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">John Doe</div>
      </div>
      <div className="flex flex-grow flex-col w-full overflow-y-scroll px-7">
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
      </div>
      <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-blue-500 min-h-20">
        <div className="flex-grow border-r-2 border-blue-500">
            <textarea ref={textareaRef} className="min-h-full w-full p-2 bg-inherit text-white outline-none resize-none" placeholder="Type a message...">
            </textarea>
        </div>
        <button onClick={handleNewMessage}>
          <div className="flex items-center px-2 py-2 bg-black flex-shrink-0 hover:bg-slate-900 hover:cursor-pointer">
              <PaperAirplaneIcon className="w-6 h-6"/>
          </div>
        </button>
      </div>
    </div>
  );
}
