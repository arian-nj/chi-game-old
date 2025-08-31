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
    ` px-12 py-8 mx-3 rounded-xl text-2xl font-semibold text-pearl backdrop-blur-md
          bg-coolgray border`
    ,
    props.selected ? 'border-orange-500 border-4 shadow-xl text-orange-800' : 'border-gray-300 text-gray-800'
  ]">
    {{ props.gameName }}
  </button>
</template>
