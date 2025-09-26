<script setup lang="ts">
import { AccountService } from '@/gen/account/v1/account_pb';
import { authTransport } from '@/lib/transport';
import { useQuery } from '@tanstack/vue-query';
import { createClient } from '@connectrpc/connect'
import { SessionService, type Time } from '@/gen/session/v1/session_pb';
import gsap from 'gsap'
import { onMounted, useTemplateRef } from 'vue';

import PlayerCard from '@/components/session/PlayerCard.vue';
import type { SessionSocket } from '@/lib/SessionWs';

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true
  }
})

const mePlayerCardRef = useTemplateRef('me-card-ref')
const oppPlayerCardRef = useTemplateRef('opp-card-ref')
// watch(() => props.isMyTurn, (newIsMyTurn) => {
//   if (newIsMyTurn) {
//     mePlayerCardRef.value?.ContinueTimer()
//     oppPlayerCardRef.value?.PauseTimer()
//   } else {
//     mePlayerCardRef.value?.PauseTimer()
//     oppPlayerCardRef.value?.ContinueTimer()
//   }
// },
//   { immediate: true }
// )

const handleTimeSync = (timeSync: Time) => {
  if (timeSync.playerId == meData.value?.account?.id && mePlayerCardRef.value) {
    mePlayerCardRef.value.totalTime = timeSync.totalTime
    mePlayerCardRef.value.spentTime = timeSync.spentTime
  }
  if (timeSync.playerId == oppData.value?.opponent?.id && oppPlayerCardRef.value) {
    oppPlayerCardRef.value.totalTime = timeSync.totalTime
    oppPlayerCardRef.value.spentTime = timeSync.spentTime
  }
}
props.sessionSocket.HandleGameTimeSyncMessage = handleTimeSync



const { isPending: meIsPending, error: meErr, data: meData } = useQuery({
  queryKey: ['me'],
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})

const { isPending: oppIsPending, error: oppErr, data: oppData } = useQuery({
  queryKey: ['session_opponent'],
  queryFn: async () => {
    const client = createClient(SessionService, authTransport)
    const data = await client.getSessionOpponent({})
    return data
  },
  staleTime: 0
})

const playerBoardRef = useTemplateRef("player-board-ref")

onMounted(() => {
  if (playerBoardRef.value) {
    gsap.from(playerBoardRef.value, {
      yPercent: -150,          // start position below
      opacity: 0,      // fade in
      duration: 1,   // animation duration
    })
  }
})

</script>

<template>
  <div ref="player-board-ref" class="
    absolute top-4 left-1/2 -translate-x-1/2
    flex items-center justify-between gap-4 p-2
    w-[90%] max-w-2xl z-50
    rounded-xl border border-slate-700
    bg-slate-900/50 shadow-2xl shadow-black/50 backdrop-blur-md
    ">

    <span v-if="meIsPending" class="flex-1 text-center text-gray-400 italic">Loading Player...</span>
    <span v-if="meErr" class="flex-1 text-center text-red-400">{{ meErr?.message }}</span>
    <PlayerCard v-else-if="meData?.account" :account="meData?.account" :isActive="true" ref="me-card-ref"
      :is-me="true" />

    <div class="flex-shrink-0 w-1.5 h-10 bg-slate-700/50 rounded-full">
    </div>

    <span v-if="oppIsPending" class="flex-1 text-center text-gray-400 italic">Finding Opponent...</span>
    <span v-if="oppErr" class="flex-1 text-center text-red-400">{{ oppErr?.message }}</span>
    <PlayerCard v-else-if="oppData?.opponent" :account="oppData?.opponent" :isActive="false" ref="opp-card-ref"
      :is-me="false" />
  </div>
</template>
