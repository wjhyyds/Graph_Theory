import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/prim',
    name: 'Prim',
    component: () => import('../views/PrimView.vue')
  },
  {
    path: '/dijkstra',
    name: 'Dijkstra',
    component: () => import('../views/DijkstraView.vue')
  },
  {
    path: '/floyd',
    name: 'Floyd',
    component: () => import('../views/FloydView.vue')
  },
  {
    path: '/bayes',
    name: 'Bayes',
    component: () => import('../views/BayesView.vue')
  },
  {
    path: '/max-matching',
    name: 'MaxMatching',
    component: () => import('../views/MaxMatchingView.vue')
  },
  {
    path: '/tsp',
    name: 'TSP',
    component: () => import('../views/TSPView.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
