<script setup lang="ts">
import XoSquare from './XoSquare.vue'

const props = defineProps({
  boardSize: Number,
  cells: Array<number>,
})

const emit = defineEmits<{
  cellSelected: [cellIndex: number]
}>()

const boardSizeClass = (props.boardSize == 5 ? "grid-cols-5" : "grid-cols-3")
let animateIndex = -1


</script>

<template>
  <div class="flex w-full items-center justify-center">
    <div :class="[`grid gap-1.5 w-4/5 lg:w-3/4  bg-gray-500`, boardSizeClass]">
      <XoSquare v-for="cell, cellIndex in props.cells" :value="cell" :animate="animateIndex == cellIndex"
        @cell_clicked="emit('cellSelected', cellIndex)" />
    </div>
  </div>
</template>
