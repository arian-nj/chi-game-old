import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import { GetJwtToken, SetJwtToken } from '@/lib/auth'
import { rawTransport } from '@/lib/transport'
import { AuthService } from '@/gen/auth/v1/auth_pb'
import { createClient } from "@connectrpc/connect"

import type { Router } from "vue-router"

import { initDataRaw, restoreInitData } from '@telegram-apps/sdk';


export function setupRouterGuards(router: Router, IsReleaseMode: boolean) {
  router.beforeEach(async (to, _, next) => {
    // Routes that don't need authentication
    if (to.name === "login") {
      return next()
    }

    if (!IsReleaseMode) {
      // Dev mode: just check for JWT
      const token = GetJwtToken()
      if (!token) {
        return next({ name: "login" })
      }
      return next()
    }

    // Release mode: validate Telegram
    const token = GetJwtToken()
    if (token) {
      return next() // already validated
    }

    try {
      const client = createClient(AuthService, rawTransport)
      restoreInitData()
      const raw = initDataRaw()
      const data = await client.validateTelegramInitData({ initData: raw })
      SetJwtToken(data.token)
      return next()
    } catch (err) {
      console.error("Telegram auth failed:", err)
      return next({ name: "login-fail" })
    }
  })
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/LoginView.vue'),
    },
    {
      path: '/finder',
      name: 'finder',
      component: () => import('../views/FinderView.vue')
    },
    {
      path: '/session',
      name: 'session',
      component: () => import('../views/SessionView.vue')
    },
    {
      path: '/conn4',
      name: 'conn4',
      component: () => import('../components/game/conn4/Conn4Board.vue')
    },
  ],
})

export default router

