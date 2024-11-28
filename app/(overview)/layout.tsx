import ChatInbox from "../ui/chat/chatInbox";
import Sidebar from "../ui/sidebar";

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex h-full min-w-96">
        <ChatInbox />
      </div>
      {children}
    </div>
  );
}
