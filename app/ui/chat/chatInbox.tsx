import { useState, useEffect } from "react";
import { NewContactModal } from "../modal/newContactModal";
import { useDebouncedCallback } from "use-debounce";
import { ChatPreview, Contact, User } from "@/app/lib/definitions";
import { useRouter } from "next/navigation";
import { CheckIcon } from "@heroicons/react/24/outline";
import clsx from "clsx";
import Preview from "./chatPreview";
import ContactCard from "./contactCard";

enum Tab {
  Chats,
  Contacts,
}

type ChatInboxProps = {
  setCurrentChat: (chat: Contact | null) => void,
}

export default function ChatInbox({ setCurrentChat }: ChatInboxProps) {
  const router = useRouter();

  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [inbox, setInbox] = useState<ChatPreview[]>([]);
  const [currentTab, setCurrentTab] = useState<Tab>(Tab.Chats);

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
          router.push("/login");
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
        router.push("/login");
        console.log(error);
      }
    };

    const fetchInbox = async () => {
      try {
        const response = await fetch(`http://localhost:3000/chat/preview`, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        if (response.ok) {
          const chatPreviews: ChatPreview[] = await response.json();
          setInbox(chatPreviews);
        } else {
          const body = await response.json();
          console.log("error");
          console.log(body);
          router.push("/login"); 
        }
      } catch (error) {
        console.log(error);
      }
    };

    if (currentTab === Tab.Chats) {
      fetchInbox();
      console.log("fetched previews");
    } else {
      fetchContacts();
      console.log("fetched contacts");
    }
  }, [currentTab]);

  const handleSearch = useDebouncedCallback(async (query: string) => {
    if (!!query || query == "") {
      return;
    }

    const response = await fetch(`http://localhost:3000/contacts?q=${query}`, {
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
  });

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
          placeholder="Search chat and contacts"
        />
      </div>
      <div className="flex">
        <div
          className={clsx(
            currentTab === Tab.Chats ? "bg-blue-600" : "border border-blue-600",
            "flex text-lg my-2 px-2 py-1 rounded-lg mx-1 hover:bg-blue-400 cursor-pointer"
          )}
          onClick={() => setCurrentTab(Tab.Chats)}
        >
          Chats
        </div>
        <div
          className={clsx(
            currentTab === Tab.Contacts ? "bg-blue-600" : "border border-blue-600",
            "flex text-lg my-2 px-2 py-1 rounded-lg mx-1 hover:bg-blue-400 cursor-pointer"
          )}
          onClick={() => setCurrentTab(Tab.Contacts)}
        >
          Contacts
        </div>
      </div>
      <div className="flex flex-col mt-3 pr-4 max-w-full h-[700px] overflow-scroll">
        {currentTab === Tab.Contacts ? (
          contacts.map((contact) => <ContactCard contact={contact} setCurrentChat={setCurrentChat}/>)
        ) : (
          inbox.map((preview) => <Preview key={preview.userId} chatPreview={preview} />)
        )}
      </div>
    </div>
  );
}