"use client";

import { useState, useEffect } from "react";
import { NewContactModal } from "../modal/newContactModal";
import { useDebouncedCallback } from "use-debounce";
import clsx from "clsx";
import { ChatPreview, Contact, User } from "@/app/lib/definitions";
import { usePathname, useRouter } from "next/navigation";
import { CheckIcon } from "@heroicons/react/24/outline";

enum Tab {
  Chats,
  Contacts,
}

export default function ChatInbox() {
  const router = useRouter();

  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [contacts, setContacts] = useState<Contact[]>([]);
  const [chatPreviews, setChatPreviews] = useState<ChatPreview[]>([]);
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

    const fetchChatPreviews = async () => {
      try {
        const response = await fetch(`http://localhost:3000/chat/preview`, {
          method: "GET",
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        if (response.ok) {
          const chatPreviews: ChatPreview[] = await response.json();
          setChatPreviews(chatPreviews);
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
      fetchChatPreviews();
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
          <Contacts items={contacts} />
        ) : (
          <ChatPreviews items={chatPreviews} />
        )}
      </div>
    </div>
  );
}

function ChatPreviews({ items }: { items: ChatPreview[] }): JSX.Element {
  const { replace } = useRouter();

  return (
    <>
      {items.map((item) => {
        return (
          <div
            onClick={() => replace(`http://localhost:3001/chat/${item.userId}`)}
            key={Number(item.userId)}
            className="flex py-2 px-2 max-h-16 w-full hover:bg-zinc-600 rounded-md cursor-default"
          >
            <div className="flex rounded-full w-12 h-12 bg-white flex-shrink-0"></div>
            <div className="ml-3 flex-grow overflow-hidden">
              <div className="flex justify-between">
                <span className="truncate pr-2 font-semibold">
                  {item.userName}
                </span>
                <span className="block text-sm flex-shrink-0">
                  {formatDate(item.sentAt)}
                </span>
              </div>

              <div className="flex items-baseline">
                {/* <CheckIcon className="w-3.5 h-3.5 flex-shrink-0 mr-1" /> */}
                <div className="w-full truncate">{item.message}</div>
                {/* <div className="bg-blue-500 text-sm text-center rounded-full flex-shrink-0 ml-1 h-5 w-5"></div> */}
              </div>
            </div>
          </div>
        );
      })}
    </>
  );
}

function Contacts({ items }: { items: Contact[] }): JSX.Element {
  const { replace } = useRouter();

  return (
    <>
      {items.map((item) => {
        return (
          <div 
            key={item.contactId}
            onClick={() => replace(`http://localhost:3001/chat/${item.contactId}`)}
            className="flex py-2 px-2 max-h-16 w-full hover:bg-zinc-600 rounded-md cursor-default">
            <div className="flex rounded-full w-12 h-12 bg-white flex-shrink-0"></div>
            <div className="ml-3 flex-grow overflow-hidden">
              <div className="flex justify-between">
                <span className="truncate pr-2 font-semibold">
                  {item.name}
                </span>
              </div>
              <div className="flex items-baseline">
                <div className="w-full truncate ml-1">Email: {item.email}</div>
              </div>
            </div>
          </div>
        );
      })}
    </>
  );
}

function formatDate(date: Date): string {
  date = new Date();
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0"); // Months are 0-indexed
  const day = date.getDate().toString().padStart(2, "0"); // Add leading zero if necessary

  return `${year}-${month}-${day}`;
}