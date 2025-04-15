import { createContext, useContext, useEffect, useRef, useState } from "react";
import { Contact, Message, MessageType, Packet, PayloadType } from "../lib/definitions_v2"
import { decodeMessageReceive, decodeMessageSentSuccess, decodePacket, encodeMessageSend, encodePacket } from "../lib/protocol_v2";
import { useAuth } from "./AuthContext";

type IMessagingContext = {
    sendMessage: (content: string) => void;
    error: string | null;
    isConnected: boolean;
    currentChatContact: Contact | null;
    currentChatMessages: Message[];
    setCurrentChatContact: (contact: Contact | null) => void;
    isCurrentChatMessagesLoading: boolean;
}

export const MessagingContext = createContext<IMessagingContext | null>(null);

export const MessagingProvider: React.FC<React.PropsWithChildren> = ({
    children
}) => {
    const socketRef = useRef<WebSocket | null>(null);
    const { user } = useAuth();
    const [error, setError] = useState<string | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const [isCurrentChatMessagesLoading, setIsCurrentChatMessagesLoading] = useState(false);
    const storedMessagesRef = useRef(new Map<number, Message[]>());
    const [currentChatMessages, setCurrentChatMessages] = useState<Message[]>([]);
    const [currentChatContact, setCurrentChatContact] = useState<Contact | null>(null);
    const currentChatContactRef = useRef<Contact | null>(null);

    // TODO: error handling
    useEffect(() => {
        let socket = new WebSocket('ws://localhost:3000/ws');
        socketRef.current = socket;

        socket.onopen = () => {
            setIsConnected(true);
            console.log("connected to message server");
        }

        socket.onerror = () => {
            socket.close(1000);
            setIsConnected(false);
            console.log("socket error");
        };

        socket.onclose = () => {
            setIsConnected(false);
            console.log("Disconnected from WebSocket. Attempting to reconnect...");
        };

        socket.onmessage = async (event) => {
            const blob: Blob = event.data;
            const arrayBuffer = await blob.arrayBuffer();
            const bytes = new Uint8Array(arrayBuffer);
            const packet = decodePacket(bytes);
            switch (packet.payloadType) {
                case PayloadType.MSG_RECEIVE:

                    // handleMessageReceive
                    const messageReceive = decodeMessageReceive(packet.payload);
                    messageReceive.localUUID = crypto.randomUUID();

                    // TODO: implement inbox system
                    // When a message is received, adjust the inbox to display the chat on top
                    // if the message thread is already in the storage, append the messages
                    // else, just present the preview in the inbox
                    if (storedMessagesRef.current.has(messageReceive.senderId)) {
                        const currentMessages = storedMessagesRef.current.get(messageReceive.senderId) ?? [];
                        const newMessages = [...currentMessages, messageReceive];
                        storedMessagesRef.current.set(messageReceive.senderId, newMessages);
                        if (currentChatContactRef.current && currentChatContactRef.current.contactId === messageReceive.senderId) {
                            setCurrentChatMessages(newMessages);
                        }
                    }
                    break;
                case PayloadType.MSG_SENT_SUCCESS:
                    console.log("success: ", decodeMessageSentSuccess(packet.payload));
                    break;
                default:
                    console.log("not supported");
                    break;
            }
        };

        return () => {
            socket.close();
        };
    }, []);

    useEffect(() => {
        currentChatContactRef.current = currentChatContact;

        if (!currentChatContact || !user) {
            return;
        }
        setIsCurrentChatMessagesLoading(true);

        const messages = storedMessagesRef.current.get(currentChatContact.contactId);
        if (!messages) {
            (async () => {
                try {
                    const response = await fetch(
                        `http://localhost:3000/chat/${currentChatContact.contactId}`,
                        {
                            method: "GET",
                            headers: { "Content-Type": "application/json" },
                            credentials: "include",
                        }
                    );

                    if (response.ok) {
                        const fetchedMessages: Message[] = await response.json();
                        fetchedMessages.forEach((message) => {
                            if (message.senderId === user.id) {
                                message.type = MessageType.Send;
                            } else {
                                message.type = MessageType.Receive;
                            }
                            message.localUUID = crypto.randomUUID();
                        });

                        storedMessagesRef.current.set(
                            currentChatContact.contactId,
                            fetchedMessages
                        );
                        setCurrentChatMessages(fetchedMessages);
                    }
                } catch (error) {
                    console.log("error: ", error);
                }
            })();
        } else {
            setCurrentChatMessages(messages);
        }
        setIsCurrentChatMessagesLoading(false);

    }, [currentChatContact]);

    const sendMessage = (content: string) => {
        if (!currentChatContact || !user) {
            return;
        }

        const socket = socketRef.current;
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            setError("Cannot send message: Socket not connected");
            return;
        }

        const messageSend: Message = {
            type: MessageType.Send,
            senderId: user.id,
            recipientId: currentChatContact.contactId,
            sentAt: new Date(),
            content: content,
            localUUID: crypto.randomUUID()
        };

        const newMessages = [...currentChatMessages, messageSend];
        setCurrentChatMessages(newMessages);
        storedMessagesRef.current.set(currentChatContact.contactId, newMessages);

        const encodedMessage = encodeMessageSend(messageSend);
        const packet: Packet = {
            version: 1,
            payloadType: PayloadType.MSG_SEND,
            payloadLength: encodedMessage.length,
            payload: encodedMessage
        }
        const encodedPacket = encodePacket(packet);
        socket.send(encodedPacket.buffer);
    }


    return (
        <MessagingContext.Provider
            value={{
                error,
                isConnected,
                isCurrentChatMessagesLoading,
                sendMessage,
                currentChatContact,
                currentChatMessages,
                setCurrentChatContact
            }}
        >
            {children}
        </MessagingContext.Provider>
    );
};

export const useMessaging = () => {
    const context = useContext(MessagingContext);
    if (!context) {
        throw new Error('useMessage must be used within an useMessaging');
    }
    return context;
}
