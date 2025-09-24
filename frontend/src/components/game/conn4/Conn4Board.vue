<script setup lang="ts">
import { useToast } from '@/components/Toast.vue';
import { ref } from 'vue';

const { toast } = useToast()

const BOARD_WIDTH = 7
const BOARD_HEIGHT = 6
const cells = ref<number[]>()

function SetStartCells() {
  cells.value = Array(BOARD_WIDTH * BOARD_HEIGHT).fill(0)
}
SetStartCells()

let lastMove = 0
function handleClick(index: number) {
  const col = index % BOARD_WIDTH

  if (!cells.value) {
    return
  }
  let targetIndex = null
  for (let rowIndex = BOARD_HEIGHT - 1; rowIndex >= 0; rowIndex--) {
    const index = rowIndex * BOARD_WIDTH + col
    const cell = cells.value[index]
    if (cell === 0) {
      targetIndex = index
      break
    }
  }
  if (targetIndex == null) {
    toast.error("پره")
    return
  }
  DoMove(targetIndex, lastMove)
  lastMove = lastMove == 1 ? 2 : 1
  // console.log(rowIndex)

  function DoMove(index: number, move: number) {
    cells.value![index] = move
  }
}

</script>

<template>
  <div class="flex justify-center items-center h-screen bg-gray-800">
    <div class="bg-gray-300 rounded-xl grid grid-cols-7 w-[95vw] max-w-[500px] p-1">
      <div v-for="cell, index in cells" :key="index" @click="handleClick(index)" :class="[
        'rounded-full aspect-square flex items-center justify-center border-gray-600',
        { 'bg-gray-800 border-4 shadow-inner': cell === 0 },
        { 'bg-gradient-to-br from-blue-500 to-blue-700 border-2 animate-drop': cell === 1 },
        { 'bg-gradient-to-br from-red-500 to-red-700 border-2 animate-drop': cell === 2 }
      ]">
      </div>
    </div>
  </div>
</template>

<style>
@keyframes drop {
  from {
    transform: translateY(-300px);
  }

  to {
    transform: translateY(0);
  }
}

.animate-drop {
  animation: drop 0.4s ease-out;
}
</style>
