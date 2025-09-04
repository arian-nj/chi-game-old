<script setup lang="ts">
import * as XoBuff from "@/gen/xo_game/v1/xo_pb";
import { SessionSocket } from '@/lib/SessionWs';
import XoGame from './XoGame.vue';
import { useTemplateRef } from "vue";
import { create, toBinary } from "@bufbuild/protobuf";
import { SessionMessageSchema } from "@/gen/session/v1/session_pb";
import { useToast } from "@/components/Toast.vue";
import Session2PlayerCards from "@/components/session/Session2PlayerCards.vue";

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true
  }
})
const XoBoardRef = useTemplateRef('xo-board-ref')
let boardSize = 3

const { toast } = useToast()
props.sessionSocket.HandleGameMessage = (gameMessage) => {
  console.log(gameMessage)
  if (gameMessage.game.case != "xo") {
    throw Error("non xo message ended up in xo")
  }
  const payload = gameMessage.game.value.payload
  switch (payload.case) {
    case "move":
      handleMoveAction(payload.value)
      break

    case "playResponse":
      handlePlayResponse(payload.value)
      break
  }
}
const handleMoveAction = (moveData: XoBuff.Move) => {
  XoBoardRef.value?.DoMove(moveData.cellIndex, moveData.cellValue)
}

const handlePlayResponse = (playResponse: XoBuff.PlayResponse) => {
  if (playResponse.isValid) {
    if (playResponse.move) {
      XoBoardRef.value?.DoMove(playResponse.move.cellIndex, playResponse.move.cellValue)
    }
  } else {
    toast.error(playResponse.reason)
  }
}



function sendClick(i: number) {
  const newSessionMsg = create(SessionMessageSchema, {
    content: {
      case: "game",
      value: {
        game: {
          case: "xo",
          value: { payload: { case: "play", value: { cellIndex: i } } }
        }
      }
    }
  })
  const bytes = toBinary(SessionMessageSchema, newSessionMsg)
  props.sessionSocket.send(bytes)
}

</script>

<template>

  <div class="flex w-full items-center justify-center">
    <Session2PlayerCards />
    <XoGame @cell-selected="sendClick" :board-size="boardSize" ref="xo-board-ref" />
  </div>
</template>
