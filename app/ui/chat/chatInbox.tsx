import { NewContactModal } from "../modal/newContactModal";
import ChatCard from "./chatcard";

export default function ChatInbox() {
    return(
        <div className="flex h-full flex-col bg-zinc-800 pt-4 pl-4 w-96">
            <div className="flex justify-between">
                <div className="text-xl font-bold">Chats</div>
                <div className="flex pr-4">
                    <NewContactModal />
                </div>
            </div>
            <div className="flex mt-3 pr-4">
                <input className="w-full h-8 rounded-md px-2 text-black" type="text" placeholder="Search chat and contacts"/>
            </div>
            <div className="flex flex-col mt-3 pr-4 max-w-full h-[700px] overflow-scroll">
                <ChatCard />
                <ChatCard />
                <ChatCard />
                <ChatCard />
                <ChatCard />
                <ChatCard />
            </div>
        </div>
    );
}