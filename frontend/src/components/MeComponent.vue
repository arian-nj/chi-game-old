<script setup lang="ts">
import { AccountService } from '@/gen/account/v1/account_pb';
import { authTransport } from '@/lib/transport';
import { createClient } from "@connectrpc/connect";
import { useQuery } from '@tanstack/vue-query';

const { isPending, error, data } = useQuery({
  queryKey: ['me'],
  staleTime: 0,
  queryFn: async () => {
    const client = createClient(AccountService, authTransport)
    const data = await client.getMe({})
    return data
  }
})
</script>


<template>
  <span v-if="isPending">Loading...</span>
  <span v-else-if="error">{{ error.message }}</span>
  <h1 v-else>{{ data?.account?.name }}</h1>
</template>
