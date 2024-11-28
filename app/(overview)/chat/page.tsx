import { ChatBubbleLeftRightIcon } from "@heroicons/react/24/solid";

export default function Chat() {
  return (
      <div className="flex flex-col justify-center items-center w-full h-full">
        <ChatBubbleLeftRightIcon className="w-24 h-24" />
        <span className="font-semibold text-lg">this should work</span>
      </div>
  );
}