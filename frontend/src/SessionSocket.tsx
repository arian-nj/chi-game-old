
export class GameSessionSocket extends WebSocket {
	constructor(url: string, protocols?: string | string[]) {
		super(url, protocols)
		this.onopen = () => {
			console.log("WebSocket connection established");
		};

		this.onmessage = async (event) => {
			const text = await event.data.text();
			console.log("Received message:", text);
		};

		this.onclose = (event) => {
			console.warn("WebSocket closed:", {
				code: event.code,
				reason: event.reason,
				wasClean: event.wasClean,
			});
		};

		this.onerror = (event) => {
			console.error("WebSocket error:", event);
		};


	}

	SendChatMessage(messageData: string) {
		const data = {
			type: "chat",
			data: messageData
		}
		const stringData = JSON.stringify(data)
		this.send(stringData)
	}
}


