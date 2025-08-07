import { useEffect, useRef, useState } from "react";
import { GameSessionSocket, type SessionEvent } from "../../lib/SessionWs";
import { Message } from "../../models/Message";
import { GetMe } from "../../lib/auth";
import { useQuery } from "@tanstack/react-query";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";

type ChatProps = {
	socketRef: React.RefObject<GameSessionSocket>
};

export function Chat({ socketRef }: ChatProps) {
	const chatRef = useRef<HTMLDivElement>(null);
	const [showMessages, setShowMessages] = useState(false);
	const [messages, setMessages] = useState<Message[]>([]);

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

	const reciveMessage = (se: SessionEvent) => {
		setMessages(prev => [...prev, new Message(se.data, 0)]);
	}
	socketRef.current.HandleChatMessage = reciveMessage

	const sendMessage = (message: string) => {
		setMessages(prev => [...prev, new Message(message, data.id)]);
		socketRef.current.SendChatMessage(message)
	}

	return (
		<div className="flex flex-col w-full h-full justify-end font-[Rubik]">
			<div className={`flex flex-col w-full h-full overflow-hidden`}>
				<div
					className={`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all duration-1000
						${showMessages ? "opacity-100 bg-gray-800/50" : "opacity-0 "}`} ref={chatRef}
				>
					{messages.map((msg, index) => (
						<ChatBubble key={index} message={msg.text} isMe={msg.userID == data.id} />
					))}
				</div>
			</div>

			<ChatInput setShowMessage={setShowMessages} sendMessage={sendMessage} />
		</div >
	);
}
