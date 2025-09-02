<script setup lang="ts">
import XoOnline from '@/components/game/xo/XoOnline.vue';
import { useToast } from '@/components/Toast.vue';
import { SessionErrorType } from '@/gen/session/v1/session_pb';

import { GetJwtToken } from "@/lib/auth";
import { GetBaseUrl } from "@/lib/baseURL";
import { SessionSocket } from "@/lib/SessionWs";
import router from '@/router/router';
import { ref } from 'vue';

const { toast } = useToast()
const isConnected = ref(false)

const sessionAPIUrl = GetBaseUrl() + "/api/session/" + "?auth_token=" + GetJwtToken()
const sessionSocket = new SessionSocket(sessionAPIUrl)

sessionSocket.onopen = () => {
  isConnected.value = true
  console.log("game session WebSocket connection established");
}
sessionSocket.onclose = (event) => {
  console.warn("WebSocket closed:", {
    code: event.code,
    reason: event.reason,
    wasClean: event.wasClean,
  });
};

sessionSocket.onerror = (event) => {
  console.error("WebSocket error:", event);
  toast.error("Socket Error")
};

function HandleError(errType: SessionErrorType) {
  if (errType == SessionErrorType.AUTH) {
    toast.error("auth error")
    router.push('/');
  } else if (errType == SessionErrorType.NOSESSION) {
    toast.error("no game")
    router.push('/');
  }
}
sessionSocket.HandleSessionErrorMessage = HandleError

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
  </div>
</template>
