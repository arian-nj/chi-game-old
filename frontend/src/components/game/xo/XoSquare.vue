<script setup lang="ts">

const props = defineProps({
  value: Number,
  animate: Boolean
})

const emit = defineEmits<{
  cell_clicked: []
}>()

function onSquareClick() {
  emit('cell_clicked')
}

</script>

<template>
  <button @click="onSquareClick" :class="[`aspect-square w-full hover:bg-gray-200 focus:bg-gray-200
      border-0
			flex items-center justify-center
			text-6xl lg:text-4xl sm:text-3xl font-extrabold
      bg-[#14bd96]
      transition-all duration-200 ease-out`
  ]">
    <span :class="[
      props.animate ? 'animate-pop' : '',
      props.value == 1 ? 'text-gray-600' : 'text-gray-100'
    ]">
      {{ props.value == 0 ? "" : (props.value == 1 ? "X" : "O") }}
    </span>
  </button>
</template>
