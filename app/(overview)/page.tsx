import Sidebar from "@/app/ui/sidebar";
import Chats from "@/app/ui/chats";

export default function Home() {
  return (
    <div className="flex h-screen">
      <Sidebar />
      <Chats />
    </div>
  );
}
