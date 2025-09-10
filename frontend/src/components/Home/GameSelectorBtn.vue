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
    ` px-12 py-8 mx-3 rounded-xl text-3xl font-bold text-text backdrop-blur-md
         border-2`
    ,
    props.selected ? ' border-text bg-teal-100/30 shadow-xl text-orange-800' : 'bg-teal-50/10 border-accent text-gray-800'
  ]">
    {{ props.gameName }}
  </button>
</template>
