import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router/router'

import { VueQueryPlugin } from '@tanstack/vue-query'
import { GetApiUrl } from './lib/baseURL'
GetApiUrl()

const app = createApp(App)

app.use(router)

app.use(VueQueryPlugin)

app.mount('#app')
