import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/solid";

export default function Home() {
  return (
      <div className="flex flex-col justify-center items-center w-full h-full">
        <ChatBubbleLeftRightIcon className="w-24 h-24" />
        <span className="font-semibold text-lg">Start a chat</span>
      </div>
  );
}