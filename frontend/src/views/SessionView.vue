<script setup lang="ts">
import XoOnline from '@/components/game/xo/XoOnline.vue';
import Chat from '@/components/chat/Chat.vue';
import { useToast } from '@/components/Toast.vue';
import { GameType, SessionErrorType } from '@/gen/session/v1/session_pb';

import { GetJwtToken } from "@/lib/auth";
import { GetApiUrl } from "@/lib/baseURL";
import { SessionSocket } from "@/lib/SessionWs";
import router from '@/router/router';
import { ref, useTemplateRef, watch } from 'vue';
import Conn4Online from '@/components/game/conn4/Conn4Online.vue';
import { log } from 'console';

const { toast } = useToast()
const isConnected = ref(false)

const activeGame = ref<null | GameType>(null)

const sessionAPIUrl = GetApiUrl() + "/api/session/" + "?auth_token=" + GetJwtToken()
const sessionSocket = new SessionSocket(sessionAPIUrl)

const ChatRef = useTemplateRef('chat-ref')

sessionSocket.onopen = () => {
  isConnected.value = true
  console.log("game session WebSocket connection established");
}

sessionSocket.onerror = (event) => {
  console.error("WebSocket error:", event);
  toast.error("Socket Error")
};

function HandleError(errType: SessionErrorType) {
  if (errType == SessionErrorType.AUTH) {
    toast.error("auth error")
    router.push('/');
  } else if (errType == SessionErrorType.NOSESSION) {
    toast.error("تو هیچ بازی ای نیستی")
    router.push('/');
  }
}
sessionSocket.HandleSessionErrorMessage = HandleError

sessionSocket.HandleChangeGametype = (changeGameMessage) => {
  const gameType = changeGameMessage.gameType
  console.log("change game type " + gameType)
  switch (gameType) {
    case GameType.XO3X3 || GameType.CONN4:

      console.log(" jk change game type " + gameType)
      activeGame.value = gameType
      break;

    default:
      return
  }
}
watch(ChatRef, () => {
  if (ChatRef.value) {
    sessionSocket.HandleChatMessage = ChatRef.value.HandleIncomingChat
  }
})

</script>


<template>
  <div class="flex w-screen h-screen items-center justify-center bg-[#14bd96]">
    <div v-if="isConnected" class="w-auto h-full overflow-hidden relative flex items-center justify-center
      aspect-[9/16] ">
      <XoOnline v-if="activeGame === GameType.XO3X3" :session-socket="sessionSocket" />
      <Conn4Online v-else-if="activeGame === GameType.CONN4" />
      <h1 v-else>No Game</h1>

    </div>

    <div v-else class="text-4xl">
      Connecting...
    </div>

    <Chat :session-socket="sessionSocket" ref='chat-ref' />
  </div>
</template>
