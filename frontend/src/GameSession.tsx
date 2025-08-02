import { useEffect, useState } from "react";
import Chat from "./Chat";
import { GameSessionSocket } from "./SessionSocket";
import { GetJwtToken } from "./lib/auth";

export default function GameSession() {
	const [socket, setSocket] = useState<GameSessionSocket | null>(null);

	useEffect(() => {
		const token = GetJwtToken()
		const wsUrl = `/api/game_session/?auth_token=${token}`;
		const controller = new AbortController();

		const checkSessionAvailability = async () => {
			try {
				const res = await fetch(wsUrl + "&noconn=false", { signal: controller.signal });
				if (!res.ok) {
					const msg = res.status === 404
						? "Session not available"
						: "Unexpected error";
					throw new Error(`${msg}: ${res.status}`);
				}
			} catch (err: any) {
				console.error("Unable to start session:", err);
				alert("Error: " + err.message);
			}
		};

		const connectWebSocket = () => {
			const conn = new GameSessionSocket(wsUrl);
			setSocket(conn)
		};

		checkSessionAvailability().then(connectWebSocket);

		return () => {
			controller.abort();
			socket?.close()
		};
	}, []);

	return <>
		<h1>Game Session</h1>
		{socket ? <Chat sessionSocket={socket} /> : <h2>socket is null</h2>}
	</>
}

