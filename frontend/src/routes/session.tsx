import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useRef, useState, type RefObject, } from 'react';
import { useQuery } from '@tanstack/react-query';
import { GetBaseUrl } from '../lib/baseURL';
import { authBeforeLoad, GetJwtToken } from '../lib/auth';
import { Chat } from '../components/Chat/Chat';
import { SessionSocket } from '../lib/SessionWs';
import XoGame from '../components/XoGame/XOGame';

export const Route = createFileRoute('/session')({
	component: RouteComponent,
	beforeLoad: authBeforeLoad,
})

function RouteComponent() {
	const sessionAPIUrl = GetBaseUrl() + "/api/session/" + "?auth_token=" + GetJwtToken()

	const socketRef = useRef<SessionSocket | null>(null)
	const [isSocketReady, setSocketReady] = useState(false)

	const { isPending, error, isSuccess } = useQuery({
		queryKey: ['checkSession'],
		queryFn: async () => {
			const response = await fetch(sessionAPIUrl + "&noconn=true")
			const status = response.status
			if (status == 404) {
				throw new Error("session not found")
			} else if (status == 500) {
				throw new Error("server Error")
			}
			return await response.json()
		}
	})

	useEffect(() => {
		if (socketRef.current != null) { return }
		const socket = new SessionSocket(sessionAPIUrl)
		socket.onopen = () => {
			console.log("WebSocket connection established");
			setSocketReady(true)
		}
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
		<div className="w-screen h-screen overflow-hidden relative bg-green-200 flex items-center justify-center">
			<div className='items-center'>
				<XoGame socketRef={socketRef as RefObject<SessionSocket>} />
			</div>
			<Chat socketRef={socketRef as RefObject<SessionSocket>} />
		</div>
	)
}

