import { createRouter, createWebHashHistory } from 'vue-router'

import Navbar from '@/components/Shared/Navbar.vue'
import Login from '@/views/Login.vue'

import Teacher from '@/views/Teacher/Teacher.vue'
import List from '@/views/Teacher/List.vue'

import Admin from '@/views/Admin/Admin.vue'
import Allocation from '@/views/Admin/Allocation.vue'
import Rank from '@/views/Admin/Rank.vue'
import Statistics from '@/views/Admin/Statistics.vue'

import Student from '@/views/Student/Student.vue'
import Selection from '@/views/Student/Selection.vue'

import SuperadminLogin from '@/views/Superadmin/Login.vue'
import SuperadminDashboard from '@/views/Superadmin/Dashboard.vue'

function ViewBuilder(View) {
  return {
    default: View,
    navbar: Navbar
  }
}

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      components: ViewBuilder(Login)
    },
    {
      path: '/config',
      name: 'config-login',
      components: { default: SuperadminLogin }
    },
    {
      path: '/config/dashboard',
      name: 'superadmin-dashboard',
      components: ViewBuilder(SuperadminDashboard)
    },
    {
      path: '/teacher',
      name: 'teacher',
      components: ViewBuilder(Teacher),
      children: [
        {
          path: 'list',
          name: 'list',
          component: List
        }
      ]
    },
    {
      path: '/admin',
      name: 'admin',
      components: ViewBuilder(Admin),
      children: [
        {
          path: 'allocation',
          name: 'allocation',
          component: Allocation
        },
        {
          path: 'rank',
          name: 'rank',
          component: Rank
        },
        {
          path: 'statistics',
          name: 'statistics',
          component: Statistics
        }
      ]
    },
    {
      path: '/student',
      name: 'student',
      components: ViewBuilder(Student),
      children: [
        {
          path: 'selection',
          name: 'selection',
          component: Selection
        }
      ]
    }
  ]
})

export { router }
