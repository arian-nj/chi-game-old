export type SessionEvent = {
	type: string;
	data: string
}

export class GameSessionSocket extends WebSocket {
	HandleChatMessage: ((msg: SessionEvent) => void) | null = null

	constructor(url: string) {
		super(url, [])
		this.onmessage = async (event) => {
			const text = await event.data.text();
			const json_data: SessionEvent = await JSON.parse(text)

			if (json_data.type === "chat") {
				if (this.HandleChatMessage) {
					this.HandleChatMessage(json_data);
				}
			} else {
				throw new Error("invalid message type " + json_data.type)
			}

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

