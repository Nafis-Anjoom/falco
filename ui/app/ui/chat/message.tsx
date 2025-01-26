import clsx from "clsx";

type MessageProps = {
    isOutgoing: boolean;
    content: string;
}

export default function Message({ isOutgoing, content }: MessageProps) {
    return (
        <div className={clsx(
            "flex w-full mt-2",
            {"justify-end": isOutgoing}
        )}>
            <div
                className={clsx(
                    "max-w-96 bg-blue-500 text-white px-4 py-2 rounded-lg",
                    {"bg-zinc-600": isOutgoing}
                )}
                >
                    {content}
            </div>
        </div>
    );
}