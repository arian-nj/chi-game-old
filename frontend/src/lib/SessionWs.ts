export type SessionSocketEvent = {
	type: string;
	data: GameAction | string
}
export type GameAction = {
	action: string
	adata: any
}

const ChatEventType = "chat" as const;
const GameActionType = "game" as const;

export class SessionSocket extends WebSocket {
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
			} else if (json_data.type == GameActionType) {
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

	SendChatMessage(text: string) {
		const data: SessionSocketEvent = {
			type: ChatEventType,
			data: text,
		}
		const stringData = JSON.stringify(data)
		this.send(stringData)
	}
	SendGameAction(gameAction: GameAction) {
		const data: SessionSocketEvent = {
			type: GameActionType,
			data: gameAction,
		}
		const stringData = JSON.stringify(data)
		this.send(stringData)
	}
}

