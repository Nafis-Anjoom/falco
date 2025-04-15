import Sidebar from "../components/Sidebar";
import ChatInbox from "../components/chat/ChatInbox_v2";
import { useEffect } from "react";
import ChatPane from "../components/chat/Chatpane";
import clsx from "clsx";
import { useNavigate } from "react-router";
import { useAuth } from "../context/AuthContext";
import { useMessaging } from "../context/MessagingContext";

export default function HomePage() {
    const { user, isLoading } = useAuth();
    const { isConnected } = useMessaging();
    const navigate = useNavigate();

    useEffect(() => {
        if (!isLoading && !user) {
            navigate('/login');
        }
    }, [isLoading, user]);

    if (isLoading) return <div>Loading...</div>;
    if (!user) return null;

    return (
        <div className="flex h-screen">
            <Sidebar />
            <div className="flex h-full min-w-96">
                <ChatInbox />
            </div>
            <ChatPane />
            {!isConnected &&
                <div
                    className={clsx("absolute bg-red-600 px-4 py-1 rounded-xl bottom-3 left-14")}
                >
                    Disconnected. Attempting to reconnect...
                </div>
            }
        </div>
    );
}
