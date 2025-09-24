<script setup lang="ts">
import type { Account } from '@/gen/account/v1/account_pb';
import gsap from 'gsap';
import { onBeforeUnmount, onMounted, ref, useTemplateRef } from 'vue';
import { useRouter } from 'vue-router';

const props = defineProps<{
  winner: Account | undefined,
  loser: Account | undefined,
}>()

const endPanelDivRef = useTemplateRef("end-panel-div-ref")
const homeBtnRef = useTemplateRef('home-btn-ref')

onMounted(() => {
  if (endPanelDivRef.value) {
    gsap.from(endPanelDivRef.value, {
      opacity: .5,
      scaleX: .5,
      scaleY: .8,
      y: 100,
      duration: .4,
    })
  }
})

onMounted(() => {
  if (homeBtnRef.value) {
    gsap.fromTo(homeBtnRef.value,
      {
        rotate: 2
      },
      {
        rotate: -2,
        duration: 1,
        yoyo: true,
        repeat: -1,
        ease: 'power1.inOut'
      },
    )
  }
})
function goToHome() {
  const router = useRouter()
  router.push('/')
}
const remainedSessionTime = ref(30)

let interval: number | undefined

onMounted(() => {
  interval = setInterval(() => {
    if (remainedSessionTime.value > 0) {
      remainedSessionTime.value--;
    } else {
      clearInterval(interval)
    }
  }, 1000)
})

onBeforeUnmount(() => {
  clearInterval(interval)
})

</script>

<template>
  <div ref="end-panel-div-ref" class="
    absolute min-w-[80%] max-w-[90%] max-h-[90%] p-8
    bg-slate-800/70 backdrop-blur-lg
    rounded-2xl border border-slate-600 shadow-2xl shadow-black/50
    flex flex-col justify-center items-center gap-6
    text-center text-white
    overflow-y-auto
  ">
    <div v-if="props.winner" class="flex flex-col items-center">
      <h1
        class="text-5xl font-black bg-gradient-to-r from-amber-400 to-yellow-300 bg-clip-text text-transparent drop-shadow-lg">
        🏆 {{ props.winner.name }}
      </h1>
    </div>

    <div v-if="props.loser" class="flex flex-col items-center opacity-60">
      <h1 class="text-4xl font-medium text-slate-400 line-through">
        💀 {{ props.loser.name }}
      </h1>
    </div>

    <hr class="w-1/2 border-slate-600 my-2" />
    <div v-if="remainedSessionTime != 0">
      <h1 class="text-4xl font-bold text-red-200">{{ remainedSessionTime }} </h1>
      <h1 class="text-lg font-bold text-red-200" dir="auto">ثانیه دیگه چت بسته میشه</h1>
    </div>
    <div v-else>
      <h1 class="text-2xl font-bold text-red-200">چت بسته شد </h1>
    </div>
    <button class="
        p-4 rounded-full text-5xl
        bg-slate-700/50 text-white
        border border-slate-600
        active:scale-95
        shadow-xl
      " ref="home-btn-ref" @click="goToHome()">
      🏠
    </button>
  </div>
</template>
