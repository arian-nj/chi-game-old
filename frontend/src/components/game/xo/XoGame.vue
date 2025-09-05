<script setup lang="ts">
import XoSquare from './XoSquare.vue'
import { ref } from 'vue';

// let nextCell = 1
// const lastCellMoved = -1
const props = defineProps<{
  boardSize: number
}>()

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

const SetCells = (newCells: number[]) => {

  console.log(newCells)
  for (let i = 0; i < newCells.length; i++)
    cells.value[i] = newCells[i]
}

defineExpose({
  DoMove,
  SetCells,
})

</script>

<template>

  <div :class="[`grid gap-1.5 w-4/5 lg:w-3/4  bg-gray-500`, boardSizeClass]">
    <XoSquare v-for="cellValue, cellIndex in cells" :value="cellValue" :animate="animateIndex == cellIndex"
      @cell_clicked="onCellSelected(cellIndex)" />
  </div>
</template>
