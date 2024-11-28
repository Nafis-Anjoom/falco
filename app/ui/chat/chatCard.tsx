import { CheckIcon } from "@heroicons/react/24/outline";

export default function ChatCard() {
  return (
    <div className="flex py-2 px-2 max-h-16 w-full hover:bg-zinc-600 rounded-md cursor-default">
      <div className="flex rounded-full w-12 h-12 bg-white flex-shrink-0"></div>
      <div className="ml-3 flex-grow overflow-hidden">
        <div className="flex justify-between">
          <span className="truncate pr-2 font-semibold">
            John Doe
          </span>
          <span className="block text-sm flex-shrink-0">2024-11-15</span>
        </div>

        <div className="flex items-baseline">
          <CheckIcon className="w-3.5 h-3.5 flex-shrink-0" />
          <div className="w-full truncate ml-1">
            A quick brown fox jumps ove a lazy dog
          </div>
          <div className="bg-blue-500 text-sm text-center rounded-full flex-shrink-0 ml-1 h-5 w-5">1</div>
        </div>
      </div>
    </div>
  );
}
