import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import React, { useEffect, useRef } from "react";
import { Contact, Message, MessageType } from "@/app/lib/definitions_v2";
import clsx from "clsx";

type ChatPaneProps = {
  contact: Contact;
  messages: Message[];
  sendMessage: (content: string) => void;
};

export default function ChatPane({ contact, messages, sendMessage }: ChatPaneProps) {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const chatPaneRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    // give the browser a moment to insert the new message before scrolling
    setTimeout(() => {
      if (chatPaneRef.current) {
        chatPaneRef.current.scrollTop = chatPaneRef.current.scrollHeight;
      }
    }, 0);
  }, [messages]);

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

    const content = textareaRef.current.value;
    if (!content || content === "") {
      return;
    }

    sendMessage(content);
    textareaRef.current.value = "";
  }

  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      <div className="flex flex-shrink-0 flex-grow-0 bg-zinc-800 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">{contact.name}</div>
      </div>
      <div ref={chatPaneRef} className="flex flex-grow flex-col w-full overflow-y-scroll px-7 pb-3">
        {messages.map((message) => {
          const output = (
            <div key={message.localUUID} className={clsx("flex w-full mt-2", {"justify-end": message.type === MessageType.Send})}>
                <div className={clsx( "max-w-96 bg-blue-500 text-white px-4 py-2 rounded-lg", {"bg-zinc-600": message.type === MessageType.Send})} >
                  {message.content}
                </div>
            </div>
          );
          return output;
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