import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      component: () => import('../views/AboutView.vue'),
    },
    {
      path: '/works',
      name: 'works',
      component: () => import('../views/works/WorksView.vue'),
    },
    {
      path: '/works/playground',
      name: 'playground',
      component: () => import('../views/works/PlayGround/MainView.vue'),
    },
  ],
})

export default router
