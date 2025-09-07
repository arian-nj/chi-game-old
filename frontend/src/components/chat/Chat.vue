<script setup lang="ts">
import ChatInput from '@/components/chat/ChatInput.vue';
import type { ChatMessage } from '@/gen/session/v1/session_pb';
import type { SessionSocket } from '@/lib/SessionWs';
import { Message } from '@/types/Message';
import { watch } from 'vue';
import { onMounted, onUnmounted, ref } from 'vue';
import ChatBubble from './ChatBubble.vue';


import { authTransport } from '@/lib/transport';
import { AccountService } from '@/gen/account/v1/account_pb';
import { createClient } from '@connectrpc/connect'
import { useQuery } from '@tanstack/vue-query';

const { isPending: meIsPending, error: meErr, data: meData } = useQuery({
  queryKey: ['me'],
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})



const allChatMessages = ref(Array<Message>())
const showChat = ref(false)

watch(showChat, () => {
  console.log("show " + showChat)
})

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true,
  }
})

function sendMessage(msgText: string) {
  props.sessionSocket.SendChatReqMessage(msgText)
  if (meData) {
    allChatMessages.value.push(new Message(msgText, meData.value!.account!.id))
  }
}

function HandleIncomingChat(chatMsg: ChatMessage) {
  allChatMessages.value.push(new Message(chatMsg.text, chatMsg.playerId))
}

function handleOutClick() {
  showChat.value = false
}
onMounted(() => {
  document.addEventListener("mousedown", handleOutClick);
})
onUnmounted(() => {
  document.removeEventListener("mousedown", handleOutClick);
})

defineExpose({
  HandleIncomingChat
})
</script>

<template>
  <div v-if="meData" class="absolute bottom-0">

    <div class="flex flex-col w-full h-full justify-end font-[Rubik]">
      <div class="flex flex-col w-full h-full overflow-hidden">
        <div :class="[`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all duration-1000`,
          showChat ? 'opacity-100 bg-gray-800/50' : 'opacity-0']">
          <ChatBubble v-for="msg in allChatMessages" :text="msg.text" :is-me="meData!.account?.id ==
            msg.userID" />
        </div>
      </div>
    </div>
    <ChatInput @input-click="showChat = true" @submit="sendMessage" />
  </div>
</template>
