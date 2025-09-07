<script setup lang="ts">

import profilePic from '@/assets/profile.jpg'

const props = defineProps({
  isMe: { required: true, type: Object as () => boolean },
  text: { required: true, type: Object as () => string }
})

</script>

<template>

  <div :class="['flex items-end gap-2.5', props.isMe ? 'flex-row-reverse' : '']">
    <img :src="profilePic" class='w-10 h-10 rounded-full' alt="profile image" />
    <div :class="['max-w-1/2 rounded-t-xl p-3 break-words', props.isMe ? 'bg-emerald-100 rounded-l-xl'
      : 'bg-gray-100 rounded-r-xl']">
      <p dir="auto">{{ props.text }}</p>
    </div>
  </div>
</template>
