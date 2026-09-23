<template>
  <div
    class="toast-container position-fixed bottom-0 end-0 p-3"
    style="z-index: 5"
  >
    <div
      class="toast hide"
      id="logoutAlert"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <div class="toast-header">
        <strong class="me-auto" id="logoutAlertLabel">登出提示</strong>
        <button type="button" class="btn-close" aria-label="Close"></button>
      </div>
      <div class="toast-body">
        {{ content }}
        <div class="mt-2 pt-2 border-top">
          <button
            type="button"
            class="btn btn-primary btn-sm"
            @click="hideToast()"
          >
            確定
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters, mapActions } from 'vuex'
import { Toast } from 'bootstrap'

export default {
  data() {
    return {
      logoutTimers: [],
      refresh: false,
      now: false
    }
  },
  watch: {
    jwt() {
      this.showToastWarning()
    }
  },
  computed: {
    ...mapState(['jwt']),
    ...mapGetters(['expireAt', 'warningTime']),
    content() {
      const { expireAt } = this
      const expireDateString = new Date(expireAt).toLocaleTimeString()
      return `系統將於${expireDateString}登出`
    },
    toast() {
      const modalElement = document.getElementById('logoutAlert')
      if (!modalElement) return null
      return new Toast(modalElement)
    }
  },
  methods: {
    ...mapActions(['refreshJWT']),
    getCurrentDateTime() {
      return Date.now()
    },

    showToastWarning() {
      const {
        logoutTimers,
        expireAt,
        logout,
        getCurrentDateTime,
        warningTime
      } = this
      const countdown = expireAt - getCurrentDateTime()

      logoutTimers.forEach(clearTimeout)

      setTimeout(
        function() {
          if (this.toast) {
            this.toast.show()
          }
        }.bind(this),
        countdown - warningTime * 1000
      )

      logoutTimers.push(logout(countdown))
    },
    onRefreshJWT() {
      this.refreshJWT()
      if (this.toast) {
        this.toast.hide()
      }
    },
    hideToast() {
      if (this.toast) {
        this.toast.hide()
      }
    },
    logout(countdown) {
      return setTimeout(function() {
        window.location.reload(true)
      }, countdown)
    }
  }
}
</script>
