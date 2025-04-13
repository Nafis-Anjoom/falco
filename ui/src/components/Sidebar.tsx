import { UserCircleIcon, ArrowLeftStartOnRectangleIcon } from "@heroicons/react/24/outline";
import { useNavigate } from "react-router";

export default function Sidebar() {
    const navigate = useNavigate();

    async function handleLogout() {
        try {
            const response = await fetch("http://localhost:3000/logout", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
            });

            if (response.ok) {
                navigate("/login");
            }
        } catch (error) {
            console.log(error);
        }
    }

    return (
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
