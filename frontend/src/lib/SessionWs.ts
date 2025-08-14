export type SessionSocketEvent = {
	type: string;
	data: any
}
export type GameAction = {
	action: string
	adata: any
}

const ChatEventType = "chat"
const GameEventType = "game"

export class GameSessionSocket extends WebSocket {
	HandleChatMessage: ((msg: SessionSocketEvent) => void) | null = null

	HandleGameAction: ((msg: GameAction) => void) | null = null

	constructor(url: string) {
		super(url, [])
		this.onmessage = async (event) => {
			const text = await event.data.text();
			const json_data: SessionSocketEvent = await JSON.parse(text)

			if (json_data.type === ChatEventType) {
				if (this.HandleChatMessage != null) {
					this.HandleChatMessage(json_data);
				} else {
					throw new Error("no chat handler is set")
				}
			} else if (json_data.type == GameEventType) {
				if (this.HandleGameAction != null) {
					this.HandleGameAction(json_data.data)
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

