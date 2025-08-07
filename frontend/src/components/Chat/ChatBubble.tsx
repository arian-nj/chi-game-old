import { detectDirection } from '../../lib/lang'
import profilePic from '../../assets/profile.jpg'

type ChatBubbleProps = {
	message: string
	isMe: boolean
}

export function ChatBubble({ message, isMe }: ChatBubbleProps) {
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
