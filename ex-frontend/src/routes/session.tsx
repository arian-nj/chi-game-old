// import { createFileRoute } from '@tanstack/react-router'
// import { useEffect, useRef, useState, type RefObject, } from 'react';
// import { useQuery } from '@tanstack/react-query';
// import { GetBaseUrl } from '../lib/baseURL';
// import { authBeforeLoad, GetJwtToken } from '../lib/auth';
// import { Chat } from '../components/Chat/Chat';
// import { SessionSocket } from '../lib/SessionWs';
// import XoGame from '../games/XO/XoGame';
// import toast from 'react-hot-toast';
//
// import { useQuery as connectQuery } from "@connectrpc/connect-query";
// import { getMe } from '../gen/account/v1/account-AccountService_connectquery';
//
// export const Route = createFileRoute('/session')({
// 	component: RouteComponent,
// 	beforeLoad: authBeforeLoad,
// })
//
// function RouteComponent() {
//
// 	const sessionAPIUrl = GetBaseUrl() + "/api/session/" + "?auth_token=" + GetJwtToken()
//
// 	const socketRef = useRef<SessionSocket | null>(null)
// 	const [isSocketReady, setSocketReady] = useState(false)
//
// 	const { isPending: isMePending, error: meError, data: meData } = connectQuery(getMe)
// 	const { isPending, error, isSuccess } = useQuery({
// 		queryKey: ['checkSession'],
// 		queryFn: async () => {
// 			const response = await fetch(sessionAPIUrl + "&noconn=true")
// 			const status = response.status
// 			if (status == 404) {
// 				throw new Error("session not found")
// 			} else if (status == 500) {
// 				throw new Error("server Error")
// 			}
// 			return await response.json()
// 		},
// 		staleTime: 0,
// 	})
//
// 	useEffect(() => {
// 		if (socketRef.current != null) { return }
// 		const socket = new SessionSocket(sessionAPIUrl)
// 		socket.onopen = () => {
// 			console.log("game session WebSocket connection established");
// 			setSocketReady(true)
// 		}
// 		socketRef.current = socket
// 	}, [isSuccess, sessionAPIUrl])
//
// 	if (isPending) {
// 		return <h1 className="text-center py-2">Connecting...</h1>
// 	}
// 	if (error) {
// 		toast.error(error.message)
// 		return <div className='text-center '>Error: {error.message}</div>
// 	}
//
// 	if (isSocketReady == false) {
// 		return (
// 			<h1>waiting for socket to connect</h1>
// 		)
// 	}
//
// 	if (isMePending) return <h1>Pending Me ...</h1>;
// 	if (meError) return <h1>error Me {String(meError)}</h1>;
// 	if (meData.account == null) {
// 		const errText = "can't find you"
// 		toast.error(errText)
// 		return <h1>{errText}</h1>
// 	}
//
//
// 	return (
// 		<div className="w-screen h-screen overflow-hidden relative flex items-center justify-center">
// 			<div className='items-center'>
// 				<XoGame socketRef={socketRef as RefObject<SessionSocket>} />
// 			</div>
// 			<Chat socketRef={socketRef as RefObject<SessionSocket>} meAccount={meData.account} />
// 		</div>
// 	)
// }

import { createFileRoute } from '@tanstack/react-router'
import { authBeforeLoad } from '../lib/auth';
import XoGame from '../games/XO/XoGame';


export const Route = createFileRoute('/session')({
	component: RouteComponent,
	beforeLoad: authBeforeLoad,
})

function RouteComponent() {

	return (
		<div className="w-screen h-screen overflow-hidden relative flex items-center justify-center">
			<div className='items-center'>
				<XoGame />
			</div>
			{/* <Chat socketRef={socketRef as RefObject<SessionSocket>} meAccount={meData.account} /> */}
		</div>
	)
}


