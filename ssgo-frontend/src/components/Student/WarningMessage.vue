<template>
  <div
    class="position-fixed top-0 start-50 translate-middle-x p-1 w-100"
    style="z-index: 5"
    @click="clearMessages"
  >
    <div
      v-for="toast in toasts"
      :key="toast.id"
      :ref="toast.name"
      class="toast hide my-1 fs-5 w-100 bg-danger text-light "
      data-bs-autohide="false"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <div class="toast-body text-center">
        <span><font-awesome-icon icon="exclamation-triangle" />{{ ' ' }}</span>
        <span>{{ toast.message }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from 'vuex'
import { Toast } from 'bootstrap'
import _ from 'lodash'

export default {
  created () {
    this.handleWarning = _.debounce(this.handleWarning, 500)
  },
  data () {
    return {
      toasts: [
        {
          name: 'confirmedLiveToast',
          message: '家長已簽署，請按「確認選科次序」。'
        },
        {
          name: 'prioritisedLiveToast',
          message: '尚未將所有選科組合排序。'
        },
        {
          name: 'signatureLiveToast',
          message: '家長必須簽署，以確認選科組合次序。'
        }
      ]
    }
  },
  watch: {
    isPrioritised () {
      this.clearMessages()
      this.handleWarning()
    },
    isSigned () {
      this.clearMessages()
      this.handleWarning()
    }
  },
  computed: {
    ...mapState('student', ['priorities', 'isSigned', 'isConfirmed']),
    ...mapGetters(['combinations']),
    isPrioritised () {
      return this.priorities.length === this.combinations.length
    },
    prioritisedToast () {
      const { prioritisedLiveToast } = this.$refs
      return new Toast(prioritisedLiveToast)
    },
    signatureToast () {
      const { signatureLiveToast } = this.$refs
      return new Toast(signatureLiveToast)
    },
    confirmedToast () {
      const { confirmedLiveToast } = this.$refs
      return new Toast(confirmedLiveToast)
    }
  },
  methods: {
    handleWarning () {
      const { prioritisedToast, signatureToast, confirmedToast } = this
      if (!this.isPrioritised) {
        prioritisedToast.show()
        return
      }

      if (!this.isSigned) {
        signatureToast.show()
        return
      }

      if (!this.isConfirmed) {
        confirmedToast.show()
      }
    },
    clearMessages () {
      const { prioritisedToast, signatureToast, confirmedToast } = this
      const toasts = [prioritisedToast, signatureToast, confirmedToast]
      toasts.forEach(toast => {
        toast.hide()
      })
    }
  }
}
</script>
