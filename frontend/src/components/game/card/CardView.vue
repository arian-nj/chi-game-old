<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from "vue";

import { CreateGame } from '@/components/game/card/cardboard';

const gameContainer = ref<HTMLDivElement | null>(null);
let game: Phaser.Game | null = null;

onMounted(() => {
  if (gameContainer.value) {
    game = CreateGame(gameContainer.value.id);
  }
});

onBeforeUnmount(() => {
  if (game) {
    game.destroy(true);
    game = null;
  }
});
</script>

<template>
  <div ref="gameContainer" id="phaser-container" class="w-full h-full"></div>
</template>

<style scoped>
#phaser-container {
  width: 800px;
  height: 600px;
}
</style>
