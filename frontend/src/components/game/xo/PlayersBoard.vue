<script setup lang="ts">
import { AccountService } from '@/gen/account/v1/account_pb';
import { authTransport } from '@/lib/transport';
import { useQuery } from '@tanstack/vue-query';
import { createClient } from '@connectrpc/connect'
import { SessionService } from '@/gen/session/v1/session_pb';
import gsap from 'gsap'
import { onMounted, useTemplateRef } from 'vue';

import * as XoBuff from "@/gen/xo_game/v1/xo_pb";
import PlayerCard from '@/components/session/PlayerCard.vue';

// const props = defineProps<{
//   isMyTurn: boolean
// }>()
//
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

const handleTimeSync = (timeSync: XoBuff.Time) => {
  if (timeSync.playerId == meData.value?.account?.id && mePlayerCardRef.value) {
    mePlayerCardRef.value.totalTime = timeSync.totalTime
    mePlayerCardRef.value.spentTime = timeSync.spentTime
  }
  if (timeSync.playerId == oppData.value?.opponent?.id && oppPlayerCardRef.value) {
    oppPlayerCardRef.value.totalTime = timeSync.totalTime
    oppPlayerCardRef.value.spentTime = timeSync.spentTime
  }
}



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

defineExpose({
  handleTimeSync: handleTimeSync
})
</script>

<template>
  <div ref="player-board-ref" class="absolute top-0 left-1/2 -translate-x-1/2
    flex px-1 py-2
    w-[85%] min-h-[10%] z-50
    border-b-4 border-x-4 border-[#2E3228]
    bg-gray-700 rounded-b-2xl shadow-md shadow-gray-700/30
    ">

    <span v-if="meIsPending">Loading ...</span>
    <span v-if="meErr">Error {{ meErr?.message }}</span>
    <PlayerCard v-else-if="meData?.account" :account="meData?.account" :isActive="true" ref="me-card-ref" />

    <div class="flex justify-center items-center bg-amber-200 p-2 mx-2 text-sm opacity-0">
      <p>new chat message</p>
    </div>

    <span v-if="oppIsPending">Loading ...</span>
    <span v-if="oppErr">Error {{ oppErr?.message }}</span>
    <PlayerCard v-else-if="oppData?.opponent" :account="oppData?.opponent" :isActive="false" ref="opp-card-ref" />
  </div>
</template>
