import { createFileRoute } from '@tanstack/react-router'
import React, { useEffect, useRef, useState, type Dispatch, type SetStateAction } from 'react';
import profilePic from '../assets/profile.jpg'
import { useQuery } from '@tanstack/react-query';
import { GetBaseUrl } from '../lib/baseURL';
import { GetJwtToken } from '../lib/auth';

export const Route = createFileRoute('/game_session')({
	component: RouteComponent,
})



export class GameSessionSocket extends WebSocket {
	constructor(url: string) {
		super(url, [])
		this.onopen = () => {
			console.log("WebSocket connection established");
		};

		this.onmessage = async (event) => {
			const text = await event.data.text();
			console.log("Received message:", text);
		};

		this.onclose = (event) => {
			console.warn("WebSocket closed:", {
				code: event.code,
				reason: event.reason,
				wasClean: event.wasClean,
			});
		};

		this.onerror = (event) => {
			console.error("WebSocket error:", event);
		};
	}

	SendChatMessage(messageData: string) {
		const data = {
			type: "chat",
			data: messageData
		}
		const stringData = JSON.stringify(data)
		this.send(stringData)
	}
}

class Message {
	text: string
	constructor(text: string) {
		this.text = text
	}
}

function RouteComponent() {
	const sessionAPIUrl = GetBaseUrl() + "/api/game_session/" + "?auth_token=" + GetJwtToken()

	const socketRef = useRef<GameSessionSocket | null>(null)

	const { isPending, error, data } = useQuery({
		queryKey: ['checkSession'],
		queryFn: async () => {
			const response = await fetch(sessionAPIUrl + "&noconn=true")
			const status = response.status
			if (status == 200) {
				console.log("ok")
			} else if (status == 404) {
				throw new Error("session not found")
			} else if (status == 500) {
				throw new Error("server Error")
			}
			return await response.json()
		}
	})

	useEffect(() => {
		if (socketRef.current) { return }
		const socket = new GameSessionSocket(sessionAPIUrl)
		socketRef.current = socket
		return () => {
			socket.close()
		}
	}, [data])

	if (isPending) {
		return <h1 className="text-center py-2">Pending...</h1>
	}
	if (error) {
		return <div className='text-center '>Error: {error.message}</div>
	}

	return (
		<>
			<h1 className="text-center py-2">Session</h1>
			<Chat socketRef={socketRef} />
		</>
	);
}

type ChatProps = {
	socketRef: React.RefObject<GameSessionSocket | null>
	// 	showMessages: boolean;
	// 	setShowMessages: Dispatch<SetStateAction<boolean>>;
	// 	messages: Message[];
	// 	setMessages: Dispatch<SetStateAction<Message[]>>;
};

function Chat({ socketRef }: ChatProps) {
	const [showMessages, setShowMessages] = useState(false);
	const [messages, setMessages] = useState<Message[]>([
		new Message("Hello!"),
		new Message("سلام!"),
		new Message("This is a dummy chat message."),
		new Message("این یک پیام آزمایشی است."),
		new Message("Hello!"),
		new Message("سلام!"),
		new Message("This is a dummy chat message."),
		new Message("این یک پیام آزمایشی است."),
		new Message("Hello!"),
		new Message("سلام!"),
		new Message("This is a dummy chat message."),
		new Message("این یک پیام آزمایشی است."),
	]);


	const chatRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const handleClick = (event: MouseEvent) => {
			if (chatRef.current && chatRef.current.contains(event.target as Node)) {
				chatRef.current.scrollTop = chatRef.current.scrollHeight;
				setShowMessages(false);
			}
		};

		document.addEventListener("mousedown", handleClick);
		return () => document.removeEventListener("mousedown", handleClick);
	}, []);

	return (
		<div className="flex flex-col w-screen h-screen font-[Rubik] overflow-hidden">

			<div className="flex flex-col flex-grow overflow-hidden">

				<div
					className={`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all duration-1000
						${showMessages ? "opacity-100 bg-gray-600" : "opacity-0 hidden"}`}
					ref={chatRef}
				>
					{messages.map((msg, index) => (
						<ChatBubble key={index} message={msg.text} isMe={index % 2 === 1} />
					))}
				</div>
			</div>

			<ChatInput setShowMessage={setShowMessages} socketRef={socketRef} setMessages={setMessages} />
		</div>
	);
}

type ChatInputProps = {
	socketRef: React.RefObject<GameSessionSocket | null>
	setShowMessage: Dispatch<SetStateAction<boolean>>;

	setMessages: Dispatch<SetStateAction<Message[]>>;
};

function ChatInput({ setShowMessage, socketRef, setMessages }: ChatInputProps) {
	const [inputMessage, setInputText] = useState("");

	const handleInput = (e: React.FormEvent<HTMLTextAreaElement>) => {
		setInputText(e.currentTarget.value);
	};

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		if (inputMessage.trim() === "") return;
		setMessages(prev => [...prev, new Message(inputMessage)]);
		setInputText("");
	};
	const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
		if (e.key === "Enter" && !e.shiftKey) {
			e.preventDefault(); // prevent newline
			setMessages(prev => [...prev, new Message(inputMessage)]);
			socketRef.current?.SendChatMessage(inputMessage)
			setInputText("");
		}
	}

	return (
		<div className="w-full border-t border-gray-300">
			<form onSubmit={handleSubmit} >
				<label htmlFor="chat" className="sr-only">Your message</label>
				<div className="flex items-center px-3 py-2 rounded-t-lg">
					<textarea
						onKeyDown={handleKeyDown}
						dir={detectDirection(inputMessage)}
						id="chat"
						rows={1}
						value={inputMessage}
						onInput={handleInput}
						onFocus={() => setShowMessage(true)}
						className="block mx-4 p-2.5 w-full text-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 bg-gray-800 text-xl"
						placeholder="Your message..."
					></textarea>

					<button type="submit" className="inline-flex justify-center p-2 text-blue-600 rounded-full cursor-pointer hover:bg-blue-100">
						<svg className="w-7 h-7 rotate-90" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 18 20">
							<path d="m17.914 18.594-8-18a1 1 0 0 0-1.828 0l-8 18a1 1 0 0 0 1.157 1.376L8 18.281V9a1 1 0 0 1 2 0v9.281l6.758 1.689a1 1 0 0 0 1.156-1.376Z" />
						</svg>
						<span className="sr-only">Send message</span>
					</button>
				</div>
			</form>
		</div>
	);
}

type ChatBubbleProps = {
	message: string
	isMe: boolean
}

function ChatBubble({ message, isMe }: ChatBubbleProps) {
	return (
		<>
			<div className={`flex ${isMe ? "flex-row-reverse" : ""} items-end gap-2.5`} >
				<img src={profilePic}
					className='w-10 h-10 rounded-full' alt="profile image"
				/>
				<div className={`max-w-1/2 rounded-t-xl p-3 break-words
					${isMe ? "rounded-l-xl" : "rounded-r-xl"}
					${isMe ? "bg-emerald-100" : "bg-gray-100"}
					`}>
					<p dir={detectDirection(message)}>{message}</p>
				</div >
			</div >
		</>
	)
}

function detectDirection(text: string): "rtl" | "ltr" {
	const firstStrongChar = text.match(/[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC\w]/);
	if (!firstStrongChar) return "ltr"; // default
	const char = firstStrongChar[0];
	return /[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC]/.test(char) ? "rtl" : "ltr";
}
