import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import Message from "./message";
import React, { useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { MessageSend, Packet, PayloadType, encodeMessageSend, encodePacket } from "@/app/lib/protocol";
import Cookies from "js-cookie";

const socketURL = "ws://localhost:3000/ws2";


export default function ChatPane() {
  const router = useRouter();
  const websocketRef = useRef<WebSocket| null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [isConnected, setIsConnected] = useState(false);
  const [messages, setMessages] = useState<MessageSend[]>([]);

  useEffect(() => {
    websocketRef.current = new WebSocket(socketURL);
    websocketRef.current.onopen = () => {
      console.log("connected to message server");
      setIsConnected(true);
    };

    websocketRef.current.onerror = () => {
      console.log("socket error");
      router.push("/login");
    }

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
    setMessages([...messages, messageSend]);

    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MessageSend,
      payloadLength: encodedMessage.length,
      payload: encodedMessage
    }
    const encodedPacket = encodePacket(packet);
    websocketRef.current?.send(encodedPacket.buffer);
  }

  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      {/* top bar */}
      <div className="flex flex-shrink-0 flex-grow-0 bg-zinc-800 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">John Doe</div>
      </div>
      {/* top bar */}
      <div className="flex flex-grow flex-col w-full overflow-y-scroll px-7">
        {messages.map((message) => {
          return (
            <Message isOutgoing={true} content={message.content}/>
          );
        })}
      </div>
      {/* <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-zinc-800 min-h-20"> */}
      <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-zinc-800">
        <div className="flex-grow ">
            <textarea ref={textareaRef} className="w-full p-2 bg-inherit text-white outline-none resize-none" placeholder="Type a message...">
            </textarea>
        </div>
        <button onClick={handleNewMessage}>
          {/* 479 */}
          <div className="flex items-center px-2 py-2 mx-4 flex-shrink-0 rounded-xl bg-blue-500 hover:bg-blue-600 hover:cursor-pointer">
              <PaperAirplaneIcon className="w-6 h-6"/>
          </div>
        </button>
      </div>
    </div>
  );
}
