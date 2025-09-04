<script setup lang="ts">
import Session2PlayerCards from '@/components/session/Session2PlayerCards.vue';
import XoSquare from './XoSquare.vue'
import { ref } from 'vue';

// let nextCell = 1
// const lastCellMoved = -1
const props = defineProps({
  boardSize: {
    type: Object as () => number,
    required: true
  },
})

const emit = defineEmits<{
  cellSelected: [cellIndex: number]
}>()


const cells = ref<number[]>(Array(props.boardSize * props.boardSize).fill(0))

const boardSizeClass = (props.boardSize == 5 ? "grid-cols-5" : "grid-cols-3")
let animateIndex = -1

function onCellSelected(index: number) {
  console.log("cell clicked", index)
  if (cells.value[index] != 0) {
    return
  }
  emit('cellSelected', index)
}

const DoMove = (index: number, value: number) => {
  cells.value[index] = value
}

defineExpose({
  DoMove,
})

</script>

<template>

  <div :class="[`grid gap-1.5 w-4/5 lg:w-3/4  bg-gray-500`, boardSizeClass]">
    <XoSquare v-for="cell, cellIndex in cells" :value="cell" :animate="animateIndex == cellIndex"
      @cell_clicked="onCellSelected(cellIndex)" />
  </div>
</template>
