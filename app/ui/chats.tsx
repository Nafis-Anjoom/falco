import ChatCard from "./chatcard";

export default function Chats() {
    return(
        <div className="flex h-full flex-col py-4 px-4 bg-yellow-700 w-96">
            <div className="bg-green text-xl font-bold">Chats</div>
            <div className="flex mt-3">
                <input className="w-full h-8 rounded-md pl-2" type="text" placeholder="Search chat"/>
            </div>
            <div className="flex flex-col mt-3 bg-green-500 max-w-full h-[700px]">
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