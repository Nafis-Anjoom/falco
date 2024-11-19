import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import Message from "./message";

export default function ChatPane() {
  return (
    <div className="flex overflow-hidden flex-col w-full h-screen">
      {/* Top bar */}
      <div className="flex flex-shrink-0 flex-grow-0 border-b-2 border-blue-500 w-full max-h-14 p-2">
        <div className="flex rounded-full w-10 h-10 bg-white flex-shrink-0"></div>
        <div className="ml-4 font-bold text-lg">John Doe</div>
      </div>
      {/* Top bar */}

      {/* thread */}
      <div className="flex flex-grow flex-col w-full overflow-y-scroll px-7">
        {/* message */}
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
        <Message isOutgoing={false}/>
        <Message isOutgoing={true}/>
      </div>
      {/* thread */}

      {/* text area */}
      <div className="flex flex-shrink-0 flex-grow-0 w-full border-t-2 border-blue-500 min-h-20">
        <div className="flex-grow border-r-2 border-blue-500">
            <textarea className="min-h-full w-full p-2 bg-inherit text-white outline-none resize-none" placeholder="Type a message...">
            </textarea>
        </div>
        <div className="flex items-center px-2 py-2 bg-black flex-shrink-0 hover:bg-slate-900 hover:cursor-pointer">
            <PaperAirplaneIcon className="w-6 h-6"/>
        </div>
      </div>
      {/* text area */}
    </div>
  );
}
