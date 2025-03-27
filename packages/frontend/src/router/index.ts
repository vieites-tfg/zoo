import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import HomePage from '../views/HomePage.vue'
import NotFound from '../views/NotFound.vue'

const routes: Array<RouteRecordRaw> = [
    {
      path: '/',
      name: 'home',
      component: HomePage
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'notFound',
      component: NotFound
    }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

export default router
