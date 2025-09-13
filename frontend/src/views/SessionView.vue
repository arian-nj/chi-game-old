<script setup lang="ts">
import XoOnline from '@/components/game/xo/XoOnline.vue';
import Chat from '@/components/chat/Chat.vue';
import { useToast } from '@/components/Toast.vue';
import { SessionErrorType } from '@/gen/session/v1/session_pb';

import { GetJwtToken } from "@/lib/auth";
import { GetApiUrl } from "@/lib/baseURL";
import { SessionSocket } from "@/lib/SessionWs";
import router from '@/router/router';
import { ref, useTemplateRef, watch } from 'vue';

const { toast } = useToast()
const isConnected = ref(false)

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
      <XoOnline :session-socket="sessionSocket" />
    </div>

    <div v-else class="text-4xl">
      Connecting...
    </div>

    <Chat :session-socket="sessionSocket" ref='chat-ref' />
  </div>
</template>
