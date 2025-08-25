import { useEffect, useRef, useState } from "react";
import { SessionSocket } from "../../lib/SessionWs";
import { ChatBubble } from "./ChatBubble";
import { ChatInput } from "./ChatInput";
import { Message } from "../../types/Message";
import { getMe } from "../../gen/account/v1/account-AccountService_connectquery";
import { useQuery } from "@connectrpc/connect-query";
import { getChatHistory } from "../../gen/session/v1/session-SessionService_connectquery";

type ChatProps = {
	socketRef: React.RefObject<SessionSocket>
};

export function Chat({ socketRef }: ChatProps) {
	const chatRef = useRef<HTMLDivElement>(null);
	const [showMessages, setShowMessages] = useState(false);
	const [messages, setMessages] = useState<Message[]>([]);

	const { isPending: isMePending, error: meError, data: meData } = useQuery(getMe)
	const { error: chatError, data: chatHistoryData } = useQuery(getChatHistory)

	useEffect(() => {
		if (chatError) {
			console.error("chat history error", chatError)
		}
		if (chatHistoryData?.messages) {
			setMessages(
				chatHistoryData.messages
					.map(
						(m) => new Message(m.text, m.playerId)
					)
			)
		}
	}, [chatHistoryData, chatError])

	// setMessages(chatHistoryData.messages.map(m => new Message(m.text, m.playerId)));


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

	// Handle click outside chat to hide
	const handleClick = (event: MouseEvent) => {
		if (chatRef.current && chatRef.current.contains(event.target as Node)) {
			setShowMessages(false);
		}
	};
	// Receive new message from socket
	socketRef.current.HandleChatMessage = (chatMessage) => {
		setMessages(prev => [...prev, new Message(chatMessage.text, chatMessage.playerId)]);
	};


	if (isMePending) return <h1>Pending Me ...</h1>;
	if (meError) return <h1>error Me {String(meError)}</h1>;

	// Send message
	const sendMessage = (message: string) => {
		setMessages(prev => [...prev, new Message(message, meData.account!.id)]);
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
								<ChatBubble key={index} message={msg.text} isMe={msg.userID === meData.account!.id} />
							))}
						</div>
					</div>
				</div>
			}
			<ChatInput setShowMessage={setShowMessages} sendMessage={sendMessage} />
		</div >
	);
}
