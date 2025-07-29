import { useEffect, useRef } from "react";

export default function GameSession() {
	const connection = useRef<WebSocket | null>(null)

	useEffect(() => {
		const socket = new WebSocket("ws://127.0.0.1:800")

		socket.onopen = (_event: any) => {
			console.log("Connection established")
		}

		socket.onmessage = (event: any) => {
			console.log("Message got", event.data)
		}

		socket.onclose = () => {
			console.log('WebSocket closed');
		};
		connection.current = socket

		return () => {
			connection.current?.close()
		}
	})
}
