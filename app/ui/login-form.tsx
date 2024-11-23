"use client";

import { lusitana } from "@/app/ui/fonts";
import {
  AtSymbolIcon,
  KeyIcon,
  ExclamationCircleIcon,
} from "@heroicons/react/24/outline";
import { ArrowRightIcon } from "@heroicons/react/20/solid";
import { Button } from "./button";
import { useRouter } from "next/navigation";
import { FormEvent } from "react";
import Link from "next/link";

interface LoginResponse {
  id: number;
  token: String;
}

export default function LoginForm() {
  const router = useRouter();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const formData = new FormData(event.currentTarget);
    const email = formData.get("email");
    const password = formData.get("password");

    console.log("email", email, "password", password);

    const response = await fetch("http://localhost:3000/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });


    if (response.ok) {
      const data = await response.json();
      localStorage.setItem("socketToken", data.token)
      router.push("/");
    } else {
      // handle errors
    }
  }

  return (
    // <form onSubmit={handleSubmit} className="space-y-3 bg-slate-900 rounded">
    <form onSubmit={handleSubmit} className="space-y-3 bg-blue-950 rounded-lg">
      <div className="flex-1 rounded-lg px-6 pb-4 pt-8">
        <h1 className={`${lusitana.className} mb-3 text-2xl`}>
          Please log in to continue.
        </h1>
        <div className="w-full">
          <div>
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm text-black outline-2 placeholder:text-gray-500"
                id="email"
                type="email"
                name="email"
                placeholder="Enter your email address"
                required
              />
              <AtSymbolIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
          </div>
          <div className="mt-4">
            <div className="relative">
              <input
                className="peer block w-full rounded-md border border-gray-200 py-[9px] pl-10 text-sm text-black outline-2 placeholder:text-gray-500"
                id="password"
                type="password"
                name="password"
                placeholder="Enter password"
                required
                minLength={6}
              />
              <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-gray-500 peer-focus:text-gray-900" />
            </div>
          </div>
        </div>
        <Button className="mt-4 w-full" aria-disabled={false}>
          Log in <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
        </Button>
        <Link href="/signup">
          <Button type="submit" className="mt-4 w-full" aria-disabled={false}>
            Sign up <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
          </Button>
        </Link>
        <div className="flex h-8 items-end space-x-1"></div>
      </div>
    </form>
  );
}
