import { Link, useNavigate } from "react-router";
import FalcoLogo from "../components/FalcoLogo";
import { FormEvent, useState } from "react";
import { ArrowRightIcon, AtSymbolIcon, ExclamationCircleIcon, KeyIcon } from "@heroicons/react/24/outline";
import { Button } from "../components/Button";
import { useAuth } from "../context/AuthContext";

export default function LoginPage() {
    return (
        <div className="min-h-screen min-w-screen flex items-center justify-center">
            <div className="relative mx-auto flex w-full max-w-[400px] flex-col space-y-2.5 p-4 md:-mt-32">
                <div className="flex h-20 w-full items-end rounded-lg bg-blue-500 p-3 md:h-36">
                    <div className="w-32 text-white md:w-36">
                        <FalcoLogo />
                    </div>
                </div>
                <LoginForm />
            </div>
        </div>
    );
}


function LoginForm() {
    const { login } = useAuth();
    const [ isLoading, setIsLoading ] = useState(false);
    const [ error, setError ] = useState<String | null>(null);
    const navigate = useNavigate();

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault();

        const formData = new FormData(event.currentTarget);
        const email = formData.get("email")?.toString() ?? "";
        const password = formData.get("password")?.toString() ?? "";

        setIsLoading(true);
        try {
            await login(email, password);
            navigate('/');
        } catch (err) {
              setError(err instanceof Error ? err.message : 'Login failed');
        } finally {
            setIsLoading(false);
        }
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-3 bg-blue-950 rounded-lg">
            <div className="flex-1 rounded-lg px-6 pb-4 pt-8">
                <h1 className={`mb-3 text-2xl`}>
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
                <Button className="mt-4 w-full" aria-disabled={isLoading} type="submit">
                    Log in <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
                </Button>
                <Link to="/signup">
                    <Button className="mt-4 w-full" aria-disabled={isLoading}>
                        Sign up <ArrowRightIcon className="ml-auto h-5 w-5 text-gray-50" />
                    </Button>
                </Link>
                <div className="flex h-8 items-end space-x-1">
                    {error && (
                        <>
                            <ExclamationCircleIcon className="h-5 w-5 text-red-500" />
                            <p className="text-sm text-red-500">{error}</p>
                        </>
                    )}
                </div>
            </div>
        </form>
    );
};
