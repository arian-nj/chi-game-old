<script setup lang="ts">
import { useToast } from '@/components/Toast.vue'
import { FinderErrorType, FinderEventSchema, FinderType } from '@/gen/finder/v1/finder_pb'
import { GetJwtToken } from '@/lib/auth'
import { GetApiUrl } from '@/lib/baseURL'
import { create, fromBinary, toBinary } from '@bufbuild/protobuf'
import { useRoute, useRouter } from 'vue-router'

import EyesAnimation from '../assets/lottie/Eyes.lottie';
import { DotLottieVue } from '@lottiefiles/dotlottie-vue'
import { DotLottie as DotLottieWeb } from '@lottiefiles/dotlottie-web'

import wasmUrl from "@lottiefiles/dotlottie-web/dist/dotlottie-player.wasm?url";
import { onMounted } from 'vue'

DotLottieWeb.setWasmUrl(wasmUrl)

onMounted(() => {
  // Prefetch Session page component
  import('../views/SessionView.vue')
})
const router = useRouter()
const route = useRoute()

let url = GetApiUrl() + "/api/match_making/ticket/?auth_token=" + GetJwtToken()
if (route.query.game) {
  url += "&game=" + route.query.game
}

const { toast } = useToast()

class FinderSocket extends WebSocket {
  constructor() {
    super(url, [])
    this.binaryType = 'arraybuffer'
  }
}
const finderSocket = new FinderSocket()

finderSocket.onopen = () => {
  console.log("WebSocket connection established");
};

finderSocket.onmessage = async (msg) => {
  const data = new Uint8Array(msg.data)
  const newFinderEvent = fromBinary(FinderEventSchema, data)
  console.log(newFinderEvent)
  if (newFinderEvent.type == FinderType.FOUND) {
    toast.success("Start")
    router.push("/session")
  } else if (newFinderEvent.type == FinderType.TIMEOUT) {
    toast.info("sorry time out")
    router.push("/")
  } else if (newFinderEvent?.errType) {
    if (newFinderEvent.errType == FinderErrorType.AUTH) {
      toast.error("invalid user")
    } else if (newFinderEvent.errType == FinderErrorType.HAS_SESSION) {
      toast.error("already have session")
    } else if (newFinderEvent.errType == FinderErrorType.HAS_TICKET) {
      toast.error("already have ticket")
    } else if (newFinderEvent.errType == FinderErrorType.INVALID_GAME) {
      toast.error("can not find this game")

    } else if (newFinderEvent.errType == FinderErrorType.UNSPECIFIED) {
      toast.error("an error happend try again")
    }
    router.push("/")
  }
}

finderSocket.onclose = (event) => {
  console.warn("WebSocket closed:", {
    code: event.code,
    reason: event.reason,
    wasClean: event.wasClean,
  });
};

finderSocket.onerror = (event) => {
  console.error("WebSocket error:", event);
  toast.error("Unxpected Error")
  router.push("/")
};

function CancelGame() {
  const cancelEvent = create(FinderEventSchema, { type: FinderType.CANCEL })
  const byte = toBinary(FinderEventSchema, cancelEvent)
  finderSocket.send(byte)
  finderSocket.close()
  router.back()
}

</script>

<template>

  <div class="flex flex-col items-center justify-center h-screen bg-neutral-400 gap-6">
    <DotLottieVue :src="EyesAnimation" loop autoplay class="w-60 sm:w-72" />

    <h1 class="text-3xl font-semibold text-gray-800 animate-pulse">
      ...دنبال حریفم
    </h1>

    <button @click="CancelGame" class='w-24 text-3xl font-bold text-gray-500'>لغو</button>
  </div>
</template>
