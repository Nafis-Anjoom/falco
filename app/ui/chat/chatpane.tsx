import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import React, { useRef } from "react";
import {
  Packet,
  PayloadType,
  encodeMessageSend,
  encodePacket,
} from "@/app/lib/protocol";
import { Chat, Contact, Message } from "@/app/lib/definitions";
import clsx from "clsx";

type ChatPaneProps = {
  contact: Contact;
  messages: Message[];
};

export default function ChatPane({ contact, messages }: ChatPaneProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const sessionUserId = BigInt(15);

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

    const messageSend: Message = {
      senderId: BigInt(15),
      recipientId: BigInt(contact.contactId),
      sentAt: new Date(),
      content: message,
    };

    // setMessages([...messages, messageSend]);

    textareaRef.current.value = "";

    const encodedMessage = encodeMessageSend(messageSend);
    const packet: Packet = {
      version: 1,
      payloadType: PayloadType.MessageSend,
      payloadLength: encodedMessage.length,
      payload: encodedMessage,
    };
    const encodedPacket = encodePacket(packet);
  }

  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      <div className="flex flex-shrink-0 flex-grow-0 bg-zinc-800 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">{contact.name}</div>
      </div>
      <div className="flex flex-grow flex-col w-full overflow-y-scroll px-7">
        {messages.map((message, index) => {
          return (
            <div key={index} className={clsx( "flex w-full mt-2", {"justify-end": message.senderId === sessionUserId})}>
                <div className={clsx( "max-w-96 bg-blue-500 text-white px-4 py-2 rounded-lg", {"bg-zinc-600": message.senderId === sessionUserId})} >
                  {message.content}
                </div>
            </div>
          );
        })}
      </div>
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
    </div>
  );
}