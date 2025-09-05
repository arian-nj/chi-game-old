<script setup lang="ts">
import type { Account } from "@/gen/account/v1/account_pb";
// import gsap from "gsap";
import { ref, watch } from "vue";

const props = defineProps<{
  isActive: Boolean,
  account: Account,
}>();

const totalTime = ref(10)
const spentTime = ref(8)

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
  const r = Math.min(255, Math.floor(255 - (progress * 2.55))); // red increases as time decreases
  const g = Math.min(255, Math.floor(progress * 2.55));         // green decreases as time decreases
  const b = 0;

  timerDiv.value.style.height = `${progress}%`
  timerDiv.value.style.backgroundColor = `rgb(${r}, ${g}, ${b})`
}

defineExpose({
  // ContinueTimer,
  // PauseTimer,
  spentTime,
  totalTime
})

</script>

<template>
  <div class="relative isolate rounded-xl overflow-hidden min-w-[25%] h-14 bg-gray-200 shadow-inner">

    <div ref="timerDiv" class="absolute bottom-0 w-full bg-blue-600"></div>

    <div class="absolute inset-0 z-10 flex items-center justify-center text-white mix-blend-difference">
      <p class="text-xl font-bold tracking-wider">{{ name }}</p>
    </div>

  </div>

</template>

// function AnimateTimer() {
// if (timerDiv.value == null) return;
//
// const remainingTime = totalTime.value - spentTime.value
// const progress = (remainingTime / totalTime.value) * 100
//
// let wasPaused = true
// if (timerTween) {
// wasPaused = timerTween?.paused()
// }
// timerTween = gsap.fromTo(timerDiv.value,
// {
// height: `${progress}%`,
// backgroundColor: startColor
// },
// {
// height: "0%",
// duration: remainingTime,
// ease: "linear",
// backgroundColor: endColor
// }
// );
// if (wasPaused) {
// timerTween.pause()
// }
// }
//
// function PauseTimer() {
// timerTween?.pause()
// }
//
// function ContinueTimer() {
// timerTween?.play()
// }
