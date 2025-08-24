import { createFileRoute, useRouter } from '@tanstack/react-router'
import { DotLottieReact } from '@lottiefiles/dotlottie-react';
import EyesAnimation from '../assets/lottie/Eyes.lottie';
import { useEffect, useRef } from 'react';
import { authBeforeLoad, GetJwtToken } from '../lib/auth';
import { GetBaseUrl } from '../lib/baseURL';
import { create, fromBinary, toBinary } from '@bufbuild/protobuf';
import { FinderEventSchema, FinderType } from '../gen/finder/v1/finder_pb';

export const Route = createFileRoute('/finder')({
	component: FinderRouteComponent,
	beforeLoad: authBeforeLoad
})

class FinderSocket extends WebSocket {
	constructor() {
		const url = GetBaseUrl() + "/api/match_making/ticket/?auth_token=" + GetJwtToken()
		super(url, [])
		this.binaryType = 'arraybuffer'
	}
}


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
		newSocket.onmessage = async (msg) => {
			const data = new Uint8Array(msg.data)
			const newFinderEvent = fromBinary(FinderEventSchema, data)
			console.log(newFinderEvent)
			if (newFinderEvent.type == FinderType.FOUND) {
				router.navigate({ to: "/session" });
			} else if (newFinderEvent.type == FinderType.TIMEOUT) {
				router.navigate({ to: "/" })
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


	}, [router])

	function cancelClicked() {


		if (socketRef.current) {
			// socketRef.current.SendMessage(FinderMessages.Cancel)
			const finderEvent = create(FinderEventSchema, { type: FinderType.CANCEL })
			const byte = toBinary(FinderEventSchema, finderEvent)
			socketRef.current.send(byte)
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

			<button
				onClick={cancelClicked}
				className='w-24 text-3xl font-bold text-gray-500' >لغو</button>
		</div>
	)
}
