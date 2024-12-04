"use client";

import { UserCircleIcon, ArrowLeftStartOnRectangleIcon } from "@heroicons/react/24/outline";
import { useRouter } from "next/navigation";

export default function Sidebar() {
    const router = useRouter();

    async function handleLogout() {
        try {
            const response = await fetch("http://localhost:3000/logout", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            });

            if (response.ok) {
                router.push("/login");
            }
        } catch(error) {
            console.log(error);
        }
    }

    return(
        <div className="flex h-full flex-col justify-between py-4 px-1.5 bg-blue-500">
            <div className="p-1">
                <UserCircleIcon className="h-8 w-8" />
            </div>
            <button onClick={() => handleLogout()} className="p-1 hover:bg-blue-400 rounded-lg">
                <ArrowLeftStartOnRectangleIcon className="h-8 w-8" />
            </button>
        </div>
    );
}