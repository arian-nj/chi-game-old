import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router, { setupRouterGuards } from './router/router'

import { VueQueryPlugin } from '@tanstack/vue-query'
import { GetApiUrl } from './lib/baseURL'
import { IsReleaseMode } from './lib/ReleaseMode'
import { init } from '@telegram-apps/sdk';

if (!IsReleaseMode) {
  init()
}
GetApiUrl()

const app = createApp(App)

app.use(router)
setupRouterGuards(router, IsReleaseMode)

app.use(VueQueryPlugin)

app.mount('#app')

