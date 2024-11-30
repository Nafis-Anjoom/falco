import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import Message from "./message";
import React, { useRef, useState } from "react";
import {
  MessageSend,
  Packet,
  PayloadType,
  encodeMessageSend,
  encodePacket,
} from "@/app/lib/protocol";
import { Contact } from "@/app/lib/definitions";

type ChatPaneProps = {
  contact: Contact;
}

export default function ChatPane({ contact }: ChatPaneProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [messages, setMessages] = useState<MessageSend[]>([]);

  function handleKeyDown(event: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (event.key === "Enter") {
      event.preventDefault();
      handleNewMessage();
    }
  }

  function handleNewMessage() {
    if (!textareaRef.current) {
      return;
    }

    const message = textareaRef.current.value;
    if (!message || message === "") {
      return;
    }

    const messageSend: MessageSend = {
      senderId: BigInt(5),
      recipientId: BigInt(1),
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
  }

  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      {/* top bar */}
      <div className="flex flex-shrink-0 flex-grow-0 bg-zinc-800 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">{contact.name}</div>
      </div>
      {/* top bar */}
      <div className="flex flex-grow flex-col w-full overflow-y-scroll px-7">
        {messages.map((message, index) => {
          return <Message key={index} isOutgoing={true} content={message.content} />;
        })}
      </div>
      {/* <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-zinc-800 min-h-20"> */}
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
          {/* 479 */}
          <div className="flex items-center px-2 py-2 mx-4 flex-shrink-0 rounded-xl bg-blue-500 hover:bg-blue-600 hover:cursor-pointer">
            <PaperAirplaneIcon className="w-6 h-6" />
          </div>
        </button>
      </div>
    </div>
  );
}