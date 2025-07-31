import { useState } from "react";
import type { GameSessionSocket } from "./SessionSocket";

type ChatProps = {
	sessionSocket: GameSessionSocket;
};

export default function Chat({ sessionSocket }: ChatProps) {

	const [message, setMessage] = useState("");

	const handleInput = (e: React.FormEvent<HTMLTextAreaElement>) => {
		setMessage(e.currentTarget.value);
	}

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		console.log("Submitted message:", message);
		sessionSocket.SendChatMessage(message)
		setMessage("");
	}

	return <div className="fixed bottom-0 left-0 w-full bg-white dark:bg-gray-800 border-t border-gray-300 dark:border-gray-700">
		<form onSubmit={handleSubmit}>
			<label htmlFor="chat" className="sr-only">Your message</label>
			<div className="flex items-center px-3 py-2 rounded-t-lg bg-gray-50 dark:bg-gray-700">
				<textarea
					dir={detectDirection(message)}
					id="chat" rows={1}
					value={message}
					onInput={handleInput}
					className="block mx-4 p-2.5 w-full text-sm text-gray-900 bg-white rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
					placeholder="Your message..."></textarea>

				<button type="submit" className="inline-flex justify-center p-2 text-blue-600 rounded-full cursor-pointer hover:bg-blue-100 dark:text-blue-500 dark:hover:bg-gray-600">
					<svg className="w-5 h-5 rotate-90 rtl:-rotate-90" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" fill="currentColor" viewBox="0 0 18 20">
						<path d="m17.914 18.594-8-18a1 1 0 0 0-1.828 0l-8 18a1 1 0 0 0 1.157 1.376L8 18.281V9a1 1 0 0 1 2 0v9.281l6.758 1.689a1 1 0 0 0 1.156-1.376Z" />
					</svg>
					<span className="sr-only">Send message</span>
				</button>
			</div>
		</form>
	</div>
}

function detectDirection(text: string): "rtl" | "ltr" {
	const firstStrongChar = text.match(/[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC\w]/);
	if (!firstStrongChar) return "ltr"; // default
	const char = firstStrongChar[0];
	return /[\u0591-\u07FF\uFB1D-\uFDFD\uFE70-\uFEFC]/.test(char) ? "rtl" : "ltr";
}


