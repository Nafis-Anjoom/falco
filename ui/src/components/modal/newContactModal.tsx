import { FormEvent, useRef } from "react";
import {
    AtSymbolIcon,
    PlusIcon,
    UserIcon,
    ExclamationCircleIcon,
    ArrowRightIcon,
    XMarkIcon
} from "@heroicons/react/24/outline";
import { useState } from "react";
import { Button } from "../../components/Button";

export function NewContactModal() {
    const modalRef = useRef<HTMLDialogElement>(null);
    const nameRef = useRef<HTMLInputElement>(null);
    const emailRef = useRef<HTMLInputElement>(null);
    const [isLoading, setIsLoading] = useState<boolean>(false);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    function handleClose() {
        if (nameRef.current) {
            nameRef.current.value = "";
        }
        if (emailRef.current) {
            emailRef.current.value = "";
        }
        setErrorMessage(null);
        modalRef.current?.close();
    }

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();
        setIsLoading(true);

        const formData = new FormData(event.currentTarget);
        const name = formData.get("name");
        const email = formData.get("email");

        const response = await fetch("http://localhost:3000/contacts", {
            method: "POST",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                name: name,
                email: email,
            }),
        });

        if (response.ok) {
            setErrorMessage(null);
            modalRef.current?.close();
        } else {
            const body = await response.json();
            setErrorMessage(body["details"]);
        }
        setIsLoading(false);
    }

    return (
        <>
            <button
                className="bg-blue-500 p-1 rounded-xl hover:bg-blue-600"
                onClick={() => modalRef.current?.showModal()}
            >
                <PlusIcon className="w-5 h-5" />
            </button>
            <dialog
                className="text-white backdrop:bg-zinc-900 backdrop:opacity-65 bg-zinc-700 rounded-lg open:fixed open:top-1/2 open:left-1/2 open:transform open:-translate-x-1/2 open:-translate-y-1/2" ref={modalRef}
            >
                <div className="w-80 h-96 p-4 flex flex-col">
                    <form
                        className="flex flex-col justify-between h-full"
                        onSubmit={handleSubmit}
                    >
                        <div>
                            <h1 className="text-center font-bold text-lg">Add contact</h1>
                            <div className="relative my-2">
                                <input
                                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm text-black outline-2 placeholder:text-gray-500"
                                    id="name"
                                    ref={nameRef}
                                    type="text"
                                    name="name"
                                    placeholder="Name"
                                    required
                                />
                                <UserIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                            </div>
                            <div className="relative my-2">
                                <input
                                    className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm text-black outline-2 placeholder:text-gray-500"
                                    id="email"
                                    ref={emailRef}
                                    type="text"
                                    name="email"
                                    placeholder="Email"
                                    required
                                />
                                <AtSymbolIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
                            </div>
                        </div>
                        <div>
                            <div className="flex h-8 items-end space-x-1">
                                {errorMessage && (
                                    <>
                                        <ExclamationCircleIcon className="h-5 w-5 text-red-500" />
                                        <p className="text-sm text-red-500">{errorMessage}</p>
                                    </>
                                )}
                            </div>
                            <Button
                                className="mt-4 w-full"
                                aria-disabled={isLoading}
                                type="submit"
                            >
                                Save{" "}
                                <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
                            </Button>

                            <Button
                                className="mt-4 w-full bg-red-600 hover:bg-red-400"
                                aria-disabled={isLoading}
                                onClick={handleClose}
                            >
                                Close{" "}
                                <XMarkIcon className="ml-auto h-5 w-5 text-gray-50" />
                            </Button>
                        </div>
                    </form>
                </div>
            </dialog>
        </>
    );
}
