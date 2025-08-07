import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useRef, useState, type RefObject, } from 'react';
import { useQuery } from '@tanstack/react-query';
import { GetBaseUrl } from '../lib/baseURL';
import { GetJwtToken } from '../lib/auth';
import { Chat } from '../components/Chat/Chat';
import { GameSessionSocket } from '../lib/SessionWs';

export const Route = createFileRoute('/game_session')({
	component: RouteComponent,
})

function RouteComponent() {
	const sessionAPIUrl = GetBaseUrl() + "/api/game_session/" + "?auth_token=" + GetJwtToken()

	const socketRef = useRef<GameSessionSocket | null>(null)
	const [isSocketReady, setSocketReady] = useState(false)

	const { isPending, error, isSuccess } = useQuery({
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
		if (socketRef.current != null) { return }
		const socket = new GameSessionSocket(sessionAPIUrl)
		socket.onopen = () => {
			console.log("WebSocket connection established");
			setSocketReady(true)
		};
		socketRef.current = socket
	}, [isSuccess])

	if (isPending) {
		return <h1 className="text-center py-2">Pending...</h1>
	}
	if (error) {
		return <div className='text-center '>Error: {error.message}</div>
	}

	if (isSocketReady == false) {
		return (
			<h1>waiting for socket to connect</h1>
		)
	}
	return (
		<div className="w-screen h-screen overflow-hidden relative bg-green-200">
			<div >
				<h1 className="text-center py-2">Session</h1>
			</div>
			<div className="absolute inset-0 z-10">
				<Chat socketRef={socketRef as RefObject<GameSessionSocket>} />
			</div>
		</div>
	);
}

