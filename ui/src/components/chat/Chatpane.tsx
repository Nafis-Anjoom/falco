import { PaperAirplaneIcon, ChatBubbleLeftRightIcon } from "@heroicons/react/24/outline";
import React, { useEffect, useRef } from "react";
import { MessageType } from "../../lib/definitions_v2";
import clsx from "clsx";
import { useMessaging } from "../../context/MessagingContext";

export default function ChatPane() {
    const textareaRef = useRef<HTMLTextAreaElement>(null);
    const chatPaneRef = useRef<HTMLDivElement>(null);
    const {
        sendMessage,
        isCurrentChatMessagesLoading,
        currentChatContact,
        currentChatMessages
    } = useMessaging();

    useEffect(() => {
        // give the browser a moment to insert the new message before scrolling
        setTimeout(() => {
            if (chatPaneRef.current) {
                chatPaneRef.current.scrollTop = chatPaneRef.current.scrollHeight;
            }
        }, 0);
    }, [currentChatMessages]);

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

    if (isCurrentChatMessagesLoading) {
        return <ChatPaneLoading />;
    }

    if (!currentChatContact) {
        return <ChatPaneSkeleton />;
    }

    return (
        <div className="flex overflow-hidden flex-col w-full h-screen">
            <div className="flex flex-shrink-0 flex-grow-0 bg-zinc-800 w-full max-h-14 p-2">
                <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
                <div className="ml-4 font-bold text-lg">{currentChatContact.name}</div>
            </div>
            <div ref={chatPaneRef} className="flex flex-grow flex-col w-full overflow-y-scroll px-7 pb-3">
                {currentChatMessages.map((message) => {
                   return (
                        <div key={message.localUUID} className={clsx("flex w-full mt-2", { "justify-end": message.type === MessageType.Send })}>
                            <div className={clsx("max-w-96 bg-blue-500 text-white px-4 py-2 rounded-lg", { "bg-zinc-600": message.type === MessageType.Send })} >
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

function ChatPaneSkeleton() {
    return (
        <div className="flex flex-col justify-center items-center w-full h-full">
            <ChatBubbleLeftRightIcon className="w-24 h-24" />
            <span className="font-semibold text-lg">Start a chat</span>
        </div>
    );
}

function ChatPaneLoading() {
    return (
        <div className="flex flex-col justify-center items-center w-full h-full">
            <ChatBubbleLeftRightIcon className="w-24 h-24" />
            <span className="font-semibold text-lg">Loading...</span>
        </div>
    );
}
