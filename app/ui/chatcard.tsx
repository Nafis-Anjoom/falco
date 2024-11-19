import { CheckIcon } from "@heroicons/react/24/outline";

export default function ChatCard() {
  return (
    <div className="flex py-2 px-2 max-h-16 bg-red-600 w-full">
      <div className="flex rounded-full w-12 h-12 bg-black flex-shrink-0"></div>
      <div className="ml-3 bg-yellow-900 flex-grow overflow-hidden">
        <div className="flex justify-between">
          <span className="truncate pr-2 font-semibold">
            John Doesd as fdsa rsefg jdsf
          </span>
          <span className="block text-sm flex-shrink-0">2024-11-15</span>
        </div>
        <div className="flex items-baseline">
          <CheckIcon className="w-3.5 h-3.5 flex-shrink-0" />
          <div className="w-full truncate ml-1">
            Chat adf asdf adsf asdf dasf dasdf is this asd fads f
          </div>
          <div className="bg-blue-500 text-sm rounded-full flex-shrink-0 ml-1 h-5 w-5">
            99
          </div>
        </div>
      </div>
    </div>
  );
}
