<script setup lang="ts">
import type { Account } from "@/gen/account/v1/account_pb";
import gsap from "gsap";
import { ref } from "vue";

const props = defineProps<{
  isActive: Boolean,
  account: Account,
}>();

defineExpose({
  AnimateTimer,
  ContinueTimer,
  PauseTimer,
})

const totalTime = ref(10)
const remaingTime = ref(8)

const startColor = "#02c432"
const endColor = "#e30232"

let name = props.account.name;
if (name.length > 10) {
  name = name.slice(0, 10) + "...";
}

const timerDiv = ref<HTMLDivElement | null>(null);

let timerTween: GSAPTween

function AnimateTimer() {
  if (timerDiv.value == null) return;
  if (timerTween) {
    timerTween.kill()
  }

  timerTween = gsap.fromTo(timerDiv.value,
    {
      height: `${(remaingTime.value / totalTime.value) * 100}%`,
      backgroundColor: startColor
    },
    {
      height: "0%",
      duration: remaingTime.value,
      ease: "linear",
      backgroundColor: endColor
    }
  );
}

function PauseTimer() {
  timerTween?.pause()
}

function ContinueTimer() {
  timerTween?.play()
}

</script>

<template>
  <div class="relative isolate rounded-xl overflow-hidden min-w-[25%] h-14 bg-gray-200 shadow-inner">

    <div ref="timerDiv" class="absolute bottom-0 w-full bg-blue-600"></div>

    <div class="absolute inset-0 z-10 flex items-center justify-center text-white mix-blend-difference">
      <p class="text-xl font-bold tracking-wider">{{ name }}</p>
    </div>

  </div>

</template>
