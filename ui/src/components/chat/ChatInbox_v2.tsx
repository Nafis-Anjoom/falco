import { useState, useEffect } from "react";
import { NewContactModal } from "../modal/newContactModal";
import { useDebouncedCallback } from "use-debounce";
import { Contact, User } from "../../lib/definitions_v2";
import { useNavigate } from "react-router";
import ContactCard from "./ContactCard";

type ChatInboxProps = {
    currentChat: Contact | null,
    setCurrentChat: (chat: Contact | null) => void,
}

export default function ChatInbox({ currentChat, setCurrentChat }: ChatInboxProps) {
    const navigate = useNavigate();

    const [currentUser, setCurrentUser] = useState<User | null>(null);
    const [contacts, setContacts] = useState<Contact[]>([]);

    useEffect(() => {
        const fetchCurrentUser = async () => {
            try {
                const response = await fetch(`http://localhost:3000/user/me`, {
                    method: "GET",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                });

                if (response.ok) {
                    const user: User = await response.json();
                    setCurrentUser(user);
                } else {
                    const body = await response.json();
                    console.log("error");
                    console.log(body);
                    navigate("/login");
                }
            } catch (error) {
                console.log(error);
            }
        };

        fetchCurrentUser();
        console.log("fetched users");
    }, []);

    useEffect(() => {
        const fetchContacts = async () => {
            try {
                const response = await fetch(`http://localhost:3000/contacts`, {
                    method: "GET",
                    headers: { "Content-Type": "application/json" },
                    credentials: "include",
                });

                if (response.ok) {
                    const contacts: Contact[] = await response.json();
                    setContacts(contacts);
                } else {
                    const body = await response.json();
                    console.log("error");
                    console.log(body);
                }
            } catch (error) {
                navigate("/login");
                console.log(error);
            }
        };

        fetchContacts();
        console.log("fetched contacts");
    }, []);

    const handleSearch = useDebouncedCallback(async (query: string) => {
        let url: string;
        if (!query || query == "") {
            url = "http://localhost:3000/contacts"
        } else {
            url = `http://localhost:3000/contacts?q=${query}`
        }

        const response = await fetch(url, {
            method: "GET",
            headers: { "Content-Type": "application/json" },
            credentials: "include",
        });

        if (response.ok) {
            const responseContacts: Contact[] = await response.json();
            setContacts(responseContacts);
        } else {
            console.log("error");
        }
    }, 300);

    return (
        <div className="flex h-full w-full flex-col bg-zinc-800 pt-4 pl-4">
            <div className="flex justify-between">
                <div className="text-xl font-bold">Hi, {currentUser?.firstName}</div>
                <div className="flex pr-4">
                    <NewContactModal />
                </div>
            </div>
            <div className="flex mt-3 pr-4">
                <input
                    className="w-full h-8 rounded-md px-2 text-black"
                    type="text"
                    placeholder="Search contacts"
                    onChange={(event) => handleSearch(event.target.value)}
                />
            </div>
            <div className="flex flex-col mt-3 pr-4 max-w-full h-[700px] overflow-scroll">
                {contacts.map((contact) => {
                    return <ContactCard
                        active={currentChat?.contactId == contact.contactId}
                        key={contact.contactId}
                        contact={contact}
                        setCurrentChat={setCurrentChat}
                    />
                })}
            </div>
        </div>
    );
}
