import { create, fromBinary, toBinary } from "@bufbuild/protobuf";
import { ChatMessageRequestSchema, SessionMessageSchema, type ChatMessage, type GameMessage } from "../gen/session/v1/session_pb";

export class SessionSocket extends WebSocket {
	HandleChatMessage: ((msg: ChatMessage) => void) | null = null

	HandleGameMessage: ((msg: GameMessage) => void) | null = null

	constructor(url: string) {
		super(url, [])
		this.binaryType = "arraybuffer"

		this.onmessage = async (event) => {
			const bytes = new Uint8Array(event.data)
			const newSessionMessage = fromBinary(SessionMessageSchema, bytes)
			if (newSessionMessage.content.case == "chat") {
				if (this.HandleChatMessage != null) {
					this.HandleChatMessage?.(newSessionMessage.content.value)
				} else {
					throw new Error("no chat handler is set")
				}
			} else if (newSessionMessage.content.case == "game") {
				if (this.HandleGameMessage != null) {
					this.HandleGameMessage(newSessionMessage.content.value)
				} else {
					console.error("no game message handler is set", newSessionMessage.content.value)
				}
			}

		}

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
	};

	SendChatMessage(text: string) {
		const chatReq = create(ChatMessageRequestSchema, { text });

		const sessionMsg = create(SessionMessageSchema, {
			content: { case: "chatReq", value: chatReq }
		});

		const bytes = toBinary(SessionMessageSchema, sessionMsg);
		this.send(bytes);
	}
	// SendGameAction(gameAction: GameAction) {
	// 	const data: SessionSocketEvent = {
	// 		type: GameActionType,
	// 		data: gameAction,
	// 	}
	// 	const stringData = JSON.stringify(data)
	// 	this.send(stringData)
	// }
	// }
}


