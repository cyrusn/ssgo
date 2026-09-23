import { createApp, h } from 'vue'
import App from '@/App.vue'
import { router } from '@/router'
import store from '@/store'
import { sync } from 'vuex-router-sync'
import 'bootstrap/dist/css/bootstrap.min.css'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

import { library } from '@fortawesome/fontawesome-svg-core'
import {
  faArrowsAlt,
  faChartBar,
  faExclamationTriangle,
  faHeart,
  faHourglass,
  faLink,
  faListAlt,
  faPrint,
  faSignInAlt,
  faSignOutAlt,
  faSitemap,
  faSyncAlt,
  faTrophy,
  faUser,
  faUserShield,
  faInfoCircle
} from '@fortawesome/free-solid-svg-icons'

sync(store, router)

library.add(
  faArrowsAlt,
  faChartBar,
  faExclamationTriangle,
  faHeart,
  faHourglass,
  faInfoCircle,
  faLink,
  faListAlt,
  faPrint,
  faSignInAlt,
  faSignOutAlt,
  faSitemap,
  faSyncAlt,
  faTrophy,
  faUser,
  faUserShield
)

// Fetch config before mounting
store.dispatch("fetchConfig")
  .catch((err) => {
    console.error("Config fetch failed (likely no active cohort):", err)
  })
  .finally(() => {
    const app = createApp({
      render() {
        return h(App)
      }
    })

    app.component('font-awesome-icon', FontAwesomeIcon)
    app.config.productionTip = false
    app.use(router)
    app.use(store)
    app.mount('#app')
  })
