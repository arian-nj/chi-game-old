<script setup>
import { useMutation } from '@tanstack/vue-query';
import { createClient } from "@connectrpc/connect";
import { ref } from 'vue';
import { SetJwtToken } from '@/lib/auth';
import { useRouter } from 'vue-router';
import { rawTransport } from '@/lib/transport';
import { DummyAuthService } from '@/gen/dummy_auth/v1/dummy_auth_pb';

const userID = ref(0)


const client = createClient(DummyAuthService, rawTransport);
const router = useRouter()
const { isPending, isError, error, isSuccess, mutate: loginMutate } = useMutation({
  mutationFn: async () => {
    const data = await client.dummyValidate({ id: userID.value });
    SetJwtToken(data?.token)
    router.back()
    return data
  }
})

function DoLogin() {
  loginMutate()
}

</script>

<template>
  <div class="about">
    <h1>This is an login page</h1>
    <h2>
      Status:
      <span v-if="isPending">Requesting...</span>
      <span v-else-if="isError">An error occurred: {{ error.message }}</span>
      <span v-else-if="isSuccess">Success</span>
      <span v-else>Enter your id</span>
    </h2>
    <input v-model="userID" class='bg-gray-200 rounded-sm' placeholder='what is your id' />
    <button type='submit' @click="DoLogin()">Send</button>

  </div>
</template>
