import { createFileRoute, useRouter } from '@tanstack/react-router'
import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import EyesAnimation from '../assets/lottie/Eyes.lottie';
import { useEffect, useRef } from 'react';
import { authBeforeLoad, GetJwtToken } from '../lib/auth';
import { GetBaseUrl } from '../lib/baseURL';

export const Route = createFileRoute('/finder')({
	component: FinderRouteComponent,
	beforeLoad: authBeforeLoad
})

class FinderSocket extends WebSocket {
	constructor() {
		const url = GetBaseUrl() + "/api/match_making/ticket/?auth_token=" + GetJwtToken()
		super(url, [])
	}

	SendMessage(message: FinderMessageValue) {
		const msg_obj = {
			type: "finder",
			data: message
		}
		const msg = JSON.stringify(msg_obj)
		console.log("sending " + msg)
		this.send(msg)
	}
}

const FinderMessages = {
	Added: "added",
	Found: "found",
	Timeout: "timeout",
	Cancel: "cancel"
}

export type FinderMessage = keyof typeof FinderMessages;

export type FinderMessageValue = (typeof FinderMessages)[keyof typeof FinderMessages];

function FinderRouteComponent() {
	const socketRef = useRef<FinderSocket | null>(null)
	const router = useRouter()

	useEffect(() => {
		if (socketRef.current != null) {
			return
		}
		const newSocket = new FinderSocket()
		socketRef.current = newSocket

		newSocket.onopen = () => {
			console.log("WebSocket connection established");
		};
		newSocket.onmessage = async (event: MessageEvent) => {
			const text = await event.data.text()
			const json_data = await JSON.parse(text)
			console.log(json_data)
			if (json_data?.type == "finder" && json_data?.data) {
				console.log("here")
				const data = json_data.data
				if (data == FinderMessages.Found) {
					router.navigate({ to: "/game_session" });
				} else if (data == FinderMessages.Timeout) {
					router.navigate({ to: "/" })
				}
			}
		}
		newSocket.onclose = (event) => {
			console.warn("WebSocket closed:", {
				code: event.code,
				reason: event.reason,
				wasClean: event.wasClean,
			});
		};

		newSocket.onerror = (event) => {
			console.error("WebSocket error:", event);
			router.navigate({ to: "/" })
		};


		// return () => newSocket.close(1000, "Client finished work")
	}, [router])

	function cancelClicked() {


		if (socketRef.current) {
			socketRef.current.SendMessage(FinderMessages.Cancel)
			socketRef.current.close()
		}
		router.navigate({ to: "/" })
	}

	return (
		<div className="flex flex-col items-center justify-center h-screen bg-gray-100 gap-6">
			<DotLottieReact
				src={EyesAnimation}
				loop
				autoplay
				className="w-60 sm:w-72"
			/>
			<h1 className="text-3xl font-semibold text-gray-800 animate-pulse">
				...دنبال حریفم
			</h1>

			{/* <div className="flex-grow"></div> */}
			<button
				onClick={cancelClicked}
				className='w-24 text-3xl font-bold text-gray-500' >لغو</button>
		</div>
	)
}
// function FinderRouteComponent() {
// 	return <div className='flex justify-center items-center h-screen'>
// 		<DotLottieReact src={EyesAnimation} loop autoplay
// 			className='w-40'
// 		/>
// 		<h1
// 			className='text-2xl font-bold'
// 		>Hello "/finder"!</h1>
// 	</div>
// }
