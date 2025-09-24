import { create, fromBinary, toBinary } from "@bufbuild/protobuf";
import { ChatMessageRequestSchema, SessionMessageSchema, type ChangeGameTypeMessage, type ChatMessage, type GameMessage, type SessionErrorType } from "../gen/session/v1/session_pb";

export class SessionSocket extends WebSocket {
  HandleChatMessage: ((chatMsg: ChatMessage) => void) | null = null

  HandleGameMessage: ((msg: GameMessage) => void) | null = null

  HandleSessionErrorMessage: ((errType: SessionErrorType) => void) | null = null

  HandleChangeGametype: ((chngGameMessage: ChangeGameTypeMessage) => void) | null = null

  constructor(url: string) {
    super(url, [])
    this.binaryType = "arraybuffer"
    this.onclose = (event) => {
      console.warn("WebSocket closed:", {
        code: event.code,
        reason: event.reason,
        wasClean: event.wasClean,
      });
    };



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
      } else if (newSessionMessage.content.case == "error") {
        if (this.HandleSessionErrorMessage != null) {
          this.HandleSessionErrorMessage(newSessionMessage.content.value)
        } else {
          console.error("no Session Error handler is set", newSessionMessage.content.value)
        }
      } else if (newSessionMessage.content.case == "gameType") {
        if (this.HandleChangeGametype != null) {
          this.HandleChangeGametype(newSessionMessage.content.value)
        }
      }

    }

  };

  SendChatReqMessage(text: string) {
    const chatReq = create(ChatMessageRequestSchema, { text });

    const sessionMsg = create(SessionMessageSchema, {
      content: { case: "chatReq", value: chatReq }
    });

    const bytes = toBinary(SessionMessageSchema, sessionMsg);
    this.send(bytes);
  }
}


