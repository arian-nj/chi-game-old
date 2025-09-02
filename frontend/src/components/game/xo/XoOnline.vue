<script setup lang="ts">
import * as XoBuff from "@/gen/xo_game/v1/xo_pb";
import { SessionSocket } from '@/lib/SessionWs';
import XoGame from './XoGame.vue';
import { ref } from "vue";
import { create, toBinary } from "@bufbuild/protobuf";
import { SessionMessageSchema } from "@/gen/session/v1/session_pb";
import { useToast } from "@/components/Toast.vue";

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true
  }
})
let boardSize = 3

let nextCell = 1
const cells = ref<number[]>(Array(boardSize * boardSize).fill(0))
const lastCellMoved = -1


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
const DoMove = (index: number, value: number) => {
  cells.value[index] = value
}
const handleMoveAction = (moveData: XoBuff.Move) => {
  DoMove(moveData.cellIndex, moveData.cellValue)
}

const handlePlayResponse = (playResponse: XoBuff.PlayResponse) => {
  if (playResponse.isValid) {
    if (playResponse.move) {
      DoMove(playResponse.move.cellIndex, playResponse.move.cellValue)
    }
  } else {
    toast.error(playResponse.reason)
  }
}



function handleClick(i: number) {
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
  // if (squares[i]) return;
  //
  // const nextSquares = squares.slice();
  // nextSquares[i] = xIsNext ? 1 : 2;
  // setSquares(nextSquares);
  // setLastPlayed(i);
}
function onCellSelected(index: number) {
  console.log("cell clicked", index)
  if (cells.value[index] != 0) {
    return
  }
  handleClick(index)
  // animateIndex = index
}


</script>

<template>
  <XoGame @cell-selected="onCellSelected" :board-size="boardSize" :cells="cells" />
</template>
