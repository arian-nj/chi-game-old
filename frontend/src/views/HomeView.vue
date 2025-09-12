<script setup lang="ts">
import { onMounted, ref } from 'vue';
import MeComponent from '../components/MeComponent.vue'
import { useRouter } from 'vue-router';
import GameSelectorBtn from '@/components/Home/GameSelectorBtn.vue';
import gsap from 'gsap'

import wasmUrl from "@lottiefiles/dotlottie-web/dist/dotlottie-player.wasm?url";


onMounted(() => {
  // prefetch
  import('../views/FinderView.vue')
  import(/* @vite-ignore */  wasmUrl)
})


const router = useRouter()


const selectedGame = ref("")
const games = ["3X3"];

const playBtnRef = ref<HTMLButtonElement | null>(null)


function handlePlayClick() {
  router.push(`/finder?game=${selectedGame.value}`)
}

onMounted(() => {
  if (playBtnRef.value) {
    gsap.from(playBtnRef.value, {
      yPercent: 100,          // start position below
      opacity: 0,      // fade in
      duration: 1,   // animation duration
    })
  }
})

</script>


<template>
  <main>
    <div class="relative h-screen w-screen bg-background text-pearl overflow-hidden">

      <div class="flex justify-center ">
        <MeComponent />
      </div>

      <div class="flex justify-center gap-4 pt-16">
        <GameSelectorBtn v-for="game in games" :game-name="game" :selected="game == selectedGame"
          @choosed="selectedGame = game" />
      </div>


      <button ref="playBtnRef" type="button" :disabled="selectedGame == ''" @click="handlePlayClick" :class="[
        `play-btn absolute bottom-0 left-0 w-full
            py-6 text-4xl font-bold rounded-none
            focus:outline-none focus:ring-4 focus:ring-pink-300
            transition-colors duration-400`
        ,
        selectedGame
          ? 'bg-secondary text-white hover:opacity-90 shadow-md'
          : 'bg-gray-300 text-gray-400 cursor-not-allowed '
      ]">
        Play
      </button>

    </div>
  </main>
</template>
