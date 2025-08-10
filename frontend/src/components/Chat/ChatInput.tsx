import { useState, type Dispatch, type SetStateAction } from "react";
import { detectDirection } from "../../lib/lang";

type ChatInputProps = {
	setShowMessage: Dispatch<SetStateAction<boolean>>;

	sendMessage: (message: string) => void
};

export function ChatInput({ setShowMessage, sendMessage }: ChatInputProps) {
	const [inputMessage, setInputText] = useState("");

	const handleInput = (e: React.FormEvent<HTMLTextAreaElement>) => {
		setInputText(e.currentTarget.value);
	};

	const send = () => {
		if (inputMessage.trim() === "") return;
		sendMessage(inputMessage)
		setInputText("");
	}

	const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault();
		send()
	};

	const handleKeyDown = (e: React.KeyboardEvent<HTMLTextAreaElement>) => {
		if ((e.key === "Enter" && !e.shiftKey) == false) { return }
		e.preventDefault();
		send()
	}

	return (
		<div className="w-full border-4">
			<form onSubmit={handleSubmit} >
				<label htmlFor="chat" className="sr-only">Your message</label>
				<div className="flex items-center px-3 py-2 rounded-t-lg">
					<textarea
						maxLength={256}
						onKeyDown={handleKeyDown}
						dir={detectDirection(inputMessage)}
						id="chat"
						rows={1}
						value={inputMessage}
						onInput={handleInput}
						onFocus={() => setShowMessage(true)}
						className="block mx-4 p-2.5 w-full text-gray-50 rounded-2xl border border-gray-300 focus:ring-blue-500 focus:border-blue-500 bg-gray-800 text-xl"
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
