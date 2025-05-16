import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Login from '../views/Login.vue'
import ManageLink from '../views/ManageLink.vue'
import ManageCategory from '../views/ManageCategory.vue'

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
      component: Login,
    },
    {
      path: '/manage',
      name: 'manage',
      component: ManageLink,
    },
    {
      path: '/manage-category',
      name: 'manage-category',
      component: ManageCategory,
    },
  ],
})

export default router
