import { Contact } from "@/app/lib/definitions_v1";
import clsx from "clsx";

type ContactCardProps = {
    contact: Contact,
    setCurrentChat: (chat: Contact | null) => void
    active?: boolean
}

export default function ContactCard({ contact, setCurrentChat, active }: ContactCardProps) {
  return (
    <div 
      onClick={() => setCurrentChat(contact)}
      className={clsx(
        "flex py-2 px-2 my-1 max-h-16 w-full hover:bg-zinc-600 rounded-md cursor-default",
        {"bg-blue-500": active}
      )}
    >
      <div className="flex rounded-full w-12 h-12 bg-white flex-shrink-0"></div>
      <div className="ml-3 flex-grow overflow-hidden">
        <div className="flex justify-between">
          <span className="truncate pr-2 font-semibold">{contact.name}</span>
        </div>
        <div className="flex items-baseline">
          <div className="w-full truncate ml-1">
            Email: {contact.email}
          </div>
        </div>
      </div>
    </div>
  );
}