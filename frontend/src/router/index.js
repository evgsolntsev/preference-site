import { createRouter, createWebHistory } from 'vue-router'
import Room from '../views/Room.vue'
import Login from '../views/Login.vue'
import Lobby from '../views/Lobby.vue'

const routes = [
  {
    path: '/room',
    name: 'room',
    component: Room
  },
  {
    path: '/lobby',
    name: 'lobby',
    component: Lobby
  },
  {
    path: '/login',
    name: 'login',
    component: Login
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
