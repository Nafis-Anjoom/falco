import clsx from "clsx";

interface MessageProps {
    isOutgoing: boolean;
}

export default function Message({ isOutgoing }: MessageProps) {
    return (
        <div className={clsx(
            "flex w-full mt-2",
            {"justify-end": isOutgoing}
        )}>
            <div className="max-w-96 w-full bg-slate-900 p-4 rounded-lg">
                Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
                Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, 
                when an unknown printer took a galley of type and scrambled it to make a type specimen book.
                It has survived not only five centuries,
            </div>
        </div>
    );
}
