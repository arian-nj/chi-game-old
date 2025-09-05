<script setup lang="ts">
import * as XoBuff from "@/gen/xo_game/v1/xo_pb";
import { SessionSocket } from '@/lib/SessionWs';
import XoGame from './XoGame.vue';
import { ref, useTemplateRef } from "vue";
import { create, toBinary } from "@bufbuild/protobuf";
import { SessionMessageSchema } from "@/gen/session/v1/session_pb";
import { useToast } from "@/components/Toast.vue";
import PlayersBoard from "./PlayersBoard.vue";

import { AccountService } from '@/gen/account/v1/account_pb';
import { authTransport } from '@/lib/transport';
import { useQuery } from '@tanstack/vue-query';
import { createClient } from '@connectrpc/connect'

const { isPending: meIsPending, error: meErr, data: meData } = useQuery({
  queryKey: ['me'],
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true
  }
})

const isMyTurn = ref(false)

const XoBoardRef = useTemplateRef('xo-board-ref')
const PlayersBoardRef = useTemplateRef('players-board')

let boardSize = ref(3)

const { toast } = useToast()
props.sessionSocket.HandleGameMessage = (gameMessage) => {
  console.log(gameMessage)
  if (gameMessage.game.case != "xo") {
    throw Error("non xo message ended up in xo")
  }
  const payload = gameMessage.game.value.payload
  switch (payload.case) {
    case "gameState":
      boardSize.value = payload.value.boardSize
      XoBoardRef.value?.SetCells(payload.value.cells)
      isMyTurn.value = meData.value?.account?.id == payload.value?.turnPlayerId
      break

    case "move":
      handleMoveAction(payload.value)
      break

    case "playResponse":
      handlePlayResponse(payload.value)
      break
    case "syncTime":
      PlayersBoardRef.value?.handleTimeSync(payload.value)
      break
  }
}
const handleMoveAction = (moveData: XoBuff.Move) => {
  XoBoardRef.value?.DoMove(moveData.cellIndex, moveData.cellValue)
  isMyTurn.value = true
}

const handlePlayResponse = (playResponse: XoBuff.PlayResponse) => {
  if (playResponse.isValid) {
    if (playResponse.move) {
      XoBoardRef.value?.DoMove(playResponse.move.cellIndex, playResponse.move.cellValue)
      isMyTurn.value = false
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

    <span v-if="meIsPending">Loading ...</span>
    <span v-if="meErr">Error {{ meErr?.message }}</span>
    <PlayersBoard v-else ref="players-board" :is-my-turn="isMyTurn" />
    <XoGame @cell-selected="sendClick" :board-size="boardSize" ref="xo-board-ref" />
  </div>
</template>
