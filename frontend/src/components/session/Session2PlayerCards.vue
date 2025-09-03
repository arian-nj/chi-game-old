<script setup lang="ts">
import { AccountService } from '@/gen/account/v1/account_pb';
import { authTransport } from '@/lib/transport';
import { useQuery } from '@tanstack/vue-query';
import { createClient } from '@connectrpc/connect'
import { SessionService } from '@/gen/session/v1/session_pb';
import PlayerCard from './PlayerCard.vue';
import gsap from 'gsap'
import { onMounted, ref } from 'vue';

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

const playerBoardRef = ref<HTMLDivElement | null>(null)

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
  <div ref="playerBoardRef" class="absolute top-0 left-1/2 -translate-x-1/2
              flex justify-between px-1 py-2
             bg-[#a7412b] border-b-4 border-x-4 border-[#2E3228]
              rounded-b-2xl shadow-md shadow-gray-700/30 w-[85%] min-h-[10%] z-50">

    <span v-if="meIsPending">Loading ...</span>
    <span v-if="meErr">Error {{ meErr?.message }}</span>
    <PlayerCard v-if="meData?.account" :account="meData.account" :isActive="true" />

    <div class="flex justify-center items-center bg-amber-200 p-2">
      <p>new chat message</p>
    </div>

    <span v-if="oppIsPending">Loading ...</span>
    <span v-if="oppErr">Error {{ oppErr?.message }}</span>
    <PlayerCard v-else-if="oppData?.opponent" :account="oppData?.opponent" :isActive="false" />


  </div>
</template>
