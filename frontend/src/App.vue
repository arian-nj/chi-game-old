<script setup lang="ts">
import { onMounted } from 'vue';
import Toast, { configureToast } from './components/Toast.vue';

onMounted(() => {
  configureToast({
    duration: 2000,
    position: 'top-center',
    pauseOnHover: false,
    limit: 2,
  });
});
</script>

<template>
  <RouterView />
  <Toast />
</template>
