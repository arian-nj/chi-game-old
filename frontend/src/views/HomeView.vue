<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import { useRouter } from 'vue-router';
import { createClient } from "@connectrpc/connect";
import GameSelectorBtn from '@/components/Home/GameSelectorBtn.vue';
import gsap from 'gsap'

import wasmUrl from "@lottiefiles/dotlottie-web/dist/dotlottie-player.wasm?url";

import { switchInlineQuery } from '@telegram-apps/sdk';
import { useQuery } from '@tanstack/vue-query';
import { SessionService } from '@/gen/session/v1/session_pb';
import { authTransport } from '@/lib/transport';
import MeComponent from '@/components/MeComponent.vue';

onMounted(() => {
  // prefetch
  import('../views/FinderView.vue')
  import(/* @vite-ignore */  wasmUrl)
})

const router = useRouter()

const games = ["3X3"];
const selectedGame = ref(games[0])

const playBtnRef = ref<HTMLButtonElement | null>(null)

function handlePlayClick() {
  if (hasSession.value) {
    router.push(`/session`)
    return
  }
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

function onPlayFriendsClick() {
  switchInlineQuery("", ["users", "groups"])
}

const { data: hasSessionData } = useQuery({
  queryKey: ['hasSession'],
  staleTime: 0,
  queryFn: async () => {
    const client = createClient(SessionService, authTransport)
    const data = await client.hasSession({})
    return data
  }
})

const hasSession = computed(() => {
  if (hasSessionData.value && hasSessionData.value.hasSession) {
    return hasSessionData.value.hasSession
  }
  return false
})

</script>

<template>
  <main>
    <div
      class="relative min-h-screen w-screen bg-gradient-to-br from-gray-950 via-gray-900 to-black text-white overflow-hidden">

      <div class="flex justify-center pt-12">
        <div class="bg-gray-900/60 p-6 rounded-2xl shadow-xl backdrop-blur-md border border-gray-700">
          <MeComponent />
        </div>
      </div>

      <div class="flex justify-center gap-6 pt-16 flex-wrap">
        <GameSelectorBtn v-for="game in games" :key="game" :game-name="game" :selected="game == selectedGame"
          @choosed="selectedGame = game" />
      </div>

      <div class="absolute bottom-0 left-0 w-full flex flex-col gap-6 items-center justify-center">
        <button class="bg-gradient-to-r from-cyan-400 to-pink-500
                  rounded-2xl px-10 py-4
                 text-2xl font-extrabold text-white tracking-wide
                 shadow-lg hover:shadow-cyan-400/30
                 hover:scale-105 active:scale-95
                 transition-all duration-300 ease-in-out" @click="onPlayFriendsClick()">
          🎮 بازی با دوستان
        </button>

        <button ref="playBtnRef" type="button" :disabled="selectedGame == ''" @click="handlePlayClick" :class="[
          `w-full py-6 text-3xl font-bold
              focus:outline-none focus:ring-4 focus:ring-pink-400/40
              transition-colors duration-400`,
          selectedGame
            ? 'bg-gradient-to-r from-emerald-400 to-green-500 text-white hover:opacity-95 shadow-xl'
            : 'bg-gray-700 text-gray-400 cursor-not-allowed'
        ]">
          {{ hasSession ?
            "ادامه بازی"
            : "🚀 شروع بازی" }}
        </button>

      </div>

    </div>
  </main>
</template>
