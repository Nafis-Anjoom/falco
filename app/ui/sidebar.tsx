import { UserCircleIcon, ArrowLeftStartOnRectangleIcon } from "@heroicons/react/24/outline";

export default function Sidebar() {
    return(
        <div className="flex h-full flex-col justify-between py-4 px-1.5 bg-red-700">
            <UserCircleIcon className="h-8 w-8" />
            <ArrowLeftStartOnRectangleIcon className="h-8 w-8" />
        </div>
    );
}