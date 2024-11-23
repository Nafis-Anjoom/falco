import { UserCircleIcon, ArrowLeftStartOnRectangleIcon } from "@heroicons/react/24/outline";

export default function Sidebar() {
    return(
        // <div className="flex h-full flex-col justify-between py-4 px-1.5 border-r-[1px] border-r-blue-500">
        <div className="flex h-full flex-col justify-between py-4 px-1.5 bg-blue-500">
            <UserCircleIcon className="h-8 w-8" />
            <ArrowLeftStartOnRectangleIcon className="h-8 w-8" />
        </div>
    );
}