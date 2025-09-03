<script setup lang="ts">

import type { Account } from "@/gen/account/v1/account_pb"

const props = defineProps<{
  isActive: Boolean,
  account: Account,
}>()
let name = props.account.name
if (name.length > 10) {
  name = name.slice(0, 10) + "..."
}
</script>
<template>

  <div class="[`rounded-xl min-w-[25%] bg-green-500 max-w-full h-full`]">
    <div class="text-white text-center bg-blue-600">
      <h3 class="text-xl font-semibold mb-1">{{ name }}</h3>
    </div>
  </div>

</template>
