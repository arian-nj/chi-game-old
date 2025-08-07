import { useEffect, useRef, useState } from "react";
import { GameSessionSocket } from "../../lib/gsWs";
import { Message } from "../../models/Message";
import { GetMe } from "../../lib/auth";
import { useQuery } from "@tanstack/react-query";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";

type ChatProps = {
	socketRef: React.RefObject<GameSessionSocket | null>
};

export function Chat({ socketRef }: ChatProps) {
	const chatRef = useRef<HTMLDivElement>(null);
	const [showMessages, setShowMessages] = useState(false);
	const [messages, setMessages] = useState<Message[]>([
		new Message("Hello!", 0),
		new Message("سلام!", 0),
		new Message("This is a dummy chat message.", 0),
		new Message("این یک پیام آزمایشی است.", 0),
		new Message("Hello!", 0),
		new Message("سلام!", 0),
		new Message("This is a dummy chat message.", 0),
		new Message("این یک پیام آزمایشی است.", 0),
		new Message("Hello!", 0),
		new Message("سلام!", 0),
		new Message("This is a dummy chat message.", 0),
		new Message("این یک پیام آزمایشی است.", 0),
	]);

	const handleClick = (event: MouseEvent) => {
		if (chatRef.current && chatRef.current.contains(event.target as Node)) {
			setShowMessages(false);
		}
	};

	useEffect(() => {
		document.addEventListener("mousedown", handleClick);
		return () => document.removeEventListener("mousedown", handleClick);
	}, []);

	useEffect(() => {
		if (showMessages && chatRef.current) {
			chatRef.current.scrollTop = chatRef.current.scrollHeight;
		}
	}, [messages, showMessages]);

	const { isPending, error, data } = useQuery({
		queryKey: ["getMe"],
		queryFn: GetMe
	})

	if (isPending) {
		return (<h1>Pending Me ...</h1>)
	}

	if (error) {
		return (<h1>error Me ${error.message}</h1>)
	}

	const meID = data?.ID
	if (!meID) return <h1>Could not get user ID</h1>;

	const sendMessage = (message: string) => {
		setMessages(prev => [...prev, new Message(message, meID)]);
		socketRef.current?.SendChatMessage(message)
	}

	return (
		<div className="flex flex-col w-full h-full justify-end font-[Rubik]">
			<div className={`flex flex-col overflow-hidden`}>
				<div
					className={`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all duration-1000
						${showMessages ? "opacity-100 bg-gray-800/50" : "opacity-0 "}`} ref={chatRef}
				>
					{messages.map((msg, index) => (
						<ChatBubble key={index} message={msg.text} isMe={msg.userID == meID} />
					))}
				</div>
			</div>

			<ChatInput setShowMessage={setShowMessages} sendMessage={sendMessage} />
		</div >
	);
}
