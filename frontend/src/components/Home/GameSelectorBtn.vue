<script setup lang="ts">
import { onMounted, ref } from 'vue';
import gsap from 'gsap'

const props = defineProps({
  gameName: String,
  selected: Boolean,
})

const emit = defineEmits(['choosed'])

const btnRef = ref<HTMLDivElement | null>()

onMounted(() => {
  if (btnRef.value) {
    gsap.from(btnRef.value, {
      scale: 0,
      opacity: 0.4,
      duration: .7,
    })
  }
})

</script>
<template>
  <button ref="btnRef" type="button" @click="emit('choosed')" :class="[
    `px-10 py-6 mx-3 rounded-2xl text-2xl font-extrabold tracking-wide
         transition-colors duration-300 ease-in-out
         border-2 backdrop-blur-md`,
    props.selected
      ? 'border-pink-500 bg-gradient-to-br from-pink-500/30 to-purple-500/30 text-white shadow-lg ring-2 ring-pink-500'
      : 'bg-gray-800/40 border-gray-600 text-gray-300 '
  ]">
    {{ props.gameName }}
  </button>
</template>
