<script setup lang="ts">
import type { Account } from '@/gen/account/v1/account_pb';
import gsap from 'gsap';
import { onMounted, useTemplateRef } from 'vue';
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
const router = useRouter()
function goToHome() {
  router.push('/')
}
</script>

<template>
  <div ref="end-panel-div-ref"
    class="absolute w-[80%] h-[50%] bg-gray-600 rounded-2xl flex flex-col justify-center items-center">

    <h1 v-if="props.winner" class="text-5xl bg-green-400 p-4 m-3 rounded-2xl border-2 border-white">🏆 {{
      props.winner.name }}</h1>

    <h1 v-if="props.loser" class="text-4xl bg-red-400 p-4 m-4 rounded-2xl border-1 border-white">💀 {{ props.loser.name
    }}</h1>

    <button class="text-7xl" ref="home-btn-ref" @click="goToHome()">🏠</button>
  </div>
</template>
