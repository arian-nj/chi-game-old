import { useEffect, useRef, useState } from "react";
import { SessionSocket, type SessionSocketEvent } from "../../lib/SessionWs";
import { fetchWithAuth, GetMe } from "../../lib/auth";
import { useQuery } from "@tanstack/react-query";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";
import { GetBaseUrl } from "../../lib/baseURL";
import { Message } from "../../types/Message";

type ChatProps = {
	socketRef: React.RefObject<SessionSocket>
};

export function Chat({ socketRef }: ChatProps) {
	const chatRef = useRef<HTMLDivElement>(null);
	const [showMessages, setShowMessages] = useState(false);
	const [messages, setMessages] = useState<Message[]>([]);

	// Get current user
	const { isPending, error, data: me } = useQuery({
		queryKey: ["getMe"],
		queryFn: GetMe
	});

	// Fetch chat history once on mount
	useEffect(() => {
		GetChatHistory()
			.then(history => {
				setMessages(history.messages.map(m => new Message(m.text, m.sender_id)));
			})
			.catch(err => {
				console.error("Failed to load chat history:", err);
			});
	}, []);

	// Handle click outside chat to hide
	const handleClick = (event: MouseEvent) => {
		if (chatRef.current && chatRef.current.contains(event.target as Node)) {
			setShowMessages(false);
		}
	};
	useEffect(() => {
		document.addEventListener("mousedown", handleClick);
		return () => document.removeEventListener("mousedown", handleClick);
	}, []);

	// Auto scroll when new messages appear
	useEffect(() => {
		if (showMessages && chatRef.current) {
			chatRef.current.scrollTop = chatRef.current.scrollHeight;
		}
	}, [messages, showMessages]);

	if (isPending) return <h1>Pending Me ...</h1>;
	if (error) return <h1>error Me {String(error)}</h1>;

	// Receive new message from socket
	const receiveMessage = (se: SessionSocketEvent) => {
		setMessages(prev => [...prev, new Message(se.data, 0)]);
	};
	socketRef.current.HandleChatMessage = receiveMessage;

	// Send message
	const sendMessage = (message: string) => {
		setMessages(prev => [...prev, new Message(message, me.id)]);
		socketRef.current.SendChatMessage(message);
	};

	return (
		<div className="absolute bottom-0" >
			{showMessages &&
				< div className="flex flex-col w-full h-full justify-end font-[Rubik]">
					<div className="flex flex-col w-full h-full overflow-hidden">
						<div
							className={`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all duration-1000
						${showMessages ? "opacity-100 bg-gray-800/50" : "opacity-0 "}`}
							ref={chatRef}
						>
							{messages.map((msg, index) => (
								<ChatBubble key={index} message={msg.text} isMe={msg.userID === me.id} />
							))}
						</div>
					</div>
				</div>
			}
			<ChatInput setShowMessage={setShowMessages} sendMessage={sendMessage} />
		</div >
	);
}

interface MessageResponse {
	id: number;
	text: string;
	sender_id: number;
}

interface ChatHistoryResponse {
	messages: Array<MessageResponse>;
}

export async function GetChatHistory() {
	const response = await fetchWithAuth(GetBaseUrl() + "/api/session/chat/history");
	if (response.status !== 200) {
		throw new Error(response.statusText);
	}
	const json_data: ChatHistoryResponse = await response.json();
	return json_data;
}
