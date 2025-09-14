<script setup lang="ts">
import type { Account } from "@/gen/account/v1/account_pb";
// import gsap from "gsap";
import { computed, ref, watch } from "vue";

const props = defineProps<{
  isActive: Boolean,
  account: Account,
  isMe: Boolean
}>();

const totalTime = ref(10)
const spentTime = ref(8)

const remainingTime = computed(() => {
  return totalTime.value - spentTime.value
})

let name = props.account.name;
if (name.length > 10) {
  name = name.slice(0, 10) + "...";
}

const timerDiv = ref<HTMLDivElement | null>(null);

// let timerTween: GSAPTween

watch(spentTime, () => {
  AnimateTimer()
})

watch(totalTime, () => {
  AnimateTimer()
})
function AnimateTimer() {
  if (timerDiv.value == null) return;

  const remainingTime = totalTime.value - spentTime.value
  const progress = (remainingTime / totalTime.value) * 100

  // calculate color from green -> red
  // progress = 100 → green (#22c55e), progress = 0 → red (#dc2626)
  const hue = progress * 1.2
  const saturation = '80%'
  const lightness = '45%'

  timerDiv.value.style.height = `${progress}%`
  timerDiv.value.style.backgroundColor = `hsl(${hue}, ${saturation}, ${lightness})`
}

defineExpose({
  spentTime,
  totalTime
})

</script>
<template>
  <div :class="[
    `relative isolate overflow-hidden flex-1 h-16
     ring-1 ring-inset ring-white/10 transition-shadow duration-300`,
    props.isMe ? 'rounded-r-lg' : 'rounded-l-lg',
    props.isActive ? 'shadow-lg shadow-cyan-500/30' : ''
  ]">

    <div ref="timerDiv" :class="[
      `absolute w-full -translate-y-1/2 top-1/2 `,
    ]"></div>

    <div class="absolute inset-0 z-10 text-white
      flex items-center justify-between px-4
      drop-shadow-lg
      ">
      <p class="text-lg font-semibold tracking-wider uppercase">{{ name }}</p>
      <p class="text-2xl font-mono font-bold">{{ remainingTime }}</p>
    </div>

  </div>
</template>
