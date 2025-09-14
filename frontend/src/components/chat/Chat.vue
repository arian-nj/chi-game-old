<script setup lang="ts">
import ChatInput from '@/components/chat/ChatInput.vue';
import { SessionService, type ChatMessage } from '@/gen/session/v1/session_pb';
import type { SessionSocket } from '@/lib/SessionWs';
import { Message } from '@/types/Message';
import { useTemplateRef, watch } from 'vue';
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

const { isPending: chatHistoryIsPending, error: chatHistoryErr, data: chatHistoryData } = useQuery({
  queryKey: ['chat-history'],
  queryFn: async () => {
    const client = createClient(SessionService, authTransport)
    const data = await client.getChatHistory({})
    return data
  },
  staleTime: 0
})

watch(chatHistoryData, () => {
  if (chatHistoryData && chatHistoryData.value) {
    const oldMessages = chatHistoryData.value.messages.map(
      (m) => new Message(m.text, m.playerId)
    )
    allChatMessages.value = oldMessages
  }
})

const allChatMessages = ref(Array<Message>())
const showChat = ref(false)
const chatInput = useTemplateRef('chat-input')
const unreadCount = ref(0); // ✨ New state for unread messages

watch(showChat, (isNowVisible) => {
  console.log("show " + isNowVisible)
  // ✨ If chat is opened, reset the counter
  if (isNowVisible) {
    unreadCount.value = 0;
  }
})

const props = defineProps({
  sessionSocket: {
    type: Object as () => SessionSocket,
    required: true,
  }
})

function sendMessage(msgText: string) {
  if (msgText == "") {
    return
  }
  props.sessionSocket.SendChatReqMessage(msgText)
  if (meData) {
    allChatMessages.value.push(new Message(msgText, meData.value!.account!.id))
  }
}

function HandleIncomingChat(chatMsg: ChatMessage) {
  allChatMessages.value.push(new Message(chatMsg.text, chatMsg.playerId))
  // ✨ Increment counter if chat is hidden
  if (!showChat.value) {
    unreadCount.value++;
  }
}

function handleOutClick(event: MouseEvent) {
  if (chatInput.value && !chatInput.value.contains(event.target as Node)) {
    showChat.value = false
  }
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
        <div :class="[`flex flex-col gap-3 px-4 py-2 overflow-y-auto flex-grow transition-all
duration-1000 rounded-xl`,
          showChat ? 'opacity-100 bg-gray-800/50' : 'opacity-0']">
          <ChatBubble v-for="msg in allChatMessages" :text="msg.text" :is-me="meData!.account?.id ==
            msg.userID" />
        </div>
      </div>
    </div>
    <div ref="chat-input" class="relative">
      <div v-if="unreadCount > 0"
        class="absolute top-0 right-0 w-5 h-5 bg-red-600 rounded-full flex items-center justify-center text-white text-xs font-bold transform translate-x-1/4 -translate-y-1/4">
        {{ unreadCount }}
      </div>

      <ChatInput @input-click="showChat = true" @submit="sendMessage" />
    </div>
  </div>
</template>
