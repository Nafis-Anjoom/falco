import { createContext, useEffect, useRef, useState } from "react";
import { Contact, Message, Packet, PayloadType } from "../lib/definitions_v2"
import { decodeMessageReceive, decodeMessageSentSuccess, decodePacket, encodeMessageSend, encodePacket } from "../lib/protocol_v2";

type IMessagingContext = {
    sendMessage: (message: Message) => void;
    error: string | null;
    isConnected: boolean;
    currentChatContact: Contact | null;
    currentChatMessages: Message[];
    setCurrentChatContact: (contact: Contact) => void;
    // isCurrentChatMessagesLoading: boolean;
    // isLoggedIn: boolean;
}

export const MessagingContext = createContext<IMessagingContext | undefined>({
    sendMessage: () => { },
    error: null,
    isConnected: false,
    currentChatContact: null,
    currentChatMessages: [],
    setCurrentChatContact: () => { },
    // isCurrentChatMessagesLoading: false,
    // isLoggedIn: false,
});

export const MessagingProvider: React.FC<React.PropsWithChildren> = ({
    children
}) => {
    const socketRef = useRef<WebSocket | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isConnected, setIsConnected] = useState(false);
    const storedMessagesRef = useRef(new Map<number, Message[]>());
    const [currentChatMessages, setCurrentChatMessages] = useState<Message[]>([]);
    const [currentChatContact, setCurrentChatContact] = useState<Contact | null>(null);

    useEffect(() => {
        const socket = new WebSocket('ws://localhost:8080');
        socketRef.current = socket;

        socket.binaryType = 'arraybuffer';

        socket.onopen = () => {
            console.log('WebSocket connected');
            setIsConnected(true);
        };

        socket.onerror = (event) => {
            console.error('WebSocket error', event);
            setError('WebSocket error');
        };

        socket.onclose = () => {
            console.log('WebSocket closed');
            setIsConnected(false);
        };

        socket.onmessage = async (event) => {
            const blob: Blob = event.data;
            const arrayBuffer = await blob.arrayBuffer();
            const bytes = new Uint8Array(arrayBuffer);
            const packet = decodePacket(bytes);
            switch (packet.payloadType) {
                case PayloadType.MSG_RECEIVE:
                    console.log('message received:', packet.payload);

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
                        if (
                            currentChatContact && currentChatContact.contactId === messageReceive.senderId
                        ) {
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

    const sendMessage = (message: Message) => {
        const socket = socketRef.current;
        if (!socket || socket.readyState !== WebSocket.OPEN) {
            setError("Cannot send message: Socket not connected");
            return;
        }

        const encodedMessage = encodeMessageSend(message);
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
