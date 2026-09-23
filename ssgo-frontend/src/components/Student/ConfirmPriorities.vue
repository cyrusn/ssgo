<template>
  <button
    type="button"
    v-if="!isConfirmed"
    class="btn me-2 my-2"
    :class="isConfirmed ? 'btn-secondary' : 'btn-danger'"
    :disabled="isConfirmed"
    @click="onClickConfirmButton"
  >
    {{ ' ' }}確定選科次序
  </button>

  <div
    class="modal fade"
    id="confirmPrioritiesModal"
    tabindex="-1"
    role="dialog"
    aria-labelledby="confirmPrioritiesModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog" role="document">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title" id="confirmPrioritiesModalLabel">
            確定選科組合次序
          </h5>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <signature-image v-if="!isConfirmed" :signature="signature" />
          <p>
            是否確定次序？一經確定後，選科組合次序一概不得更改。
          </p>
        </div>

        <div class="modal-footer">
          <button
            type="button"
            class="btn btn-secondary"
            data-bs-dismiss="modal"
          >
            取消
          </button>
          <button
            type="button"
            class="btn btn-danger"
            @click="onConfirm"
            :disabled="isConfirmed || !isPrioritised"
          >
            確定
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { Modal } from 'bootstrap'
import SignatureImage from '@/components/Shared/SignatureImage'
import { mapState, mapGetters, mapActions } from 'vuex'

export default {
  components: { SignatureImage },
  watch: {
    isPrioritised () {
      console.log(this.isPrioritised)
      if (!this.isPrioritised) {
        this.setSignature({
          isSigned: false,
          data: ''
        })
      }
    }
  },
  computed: {
    modal () {
      const modalElement = document.getElementById('confirmPrioritiesModal')
      return new Modal(modalElement)
    },
    ...mapState('student', [
      'priorities',
      'isConfirmed',
      'isSigned',
      'signature'
    ]),
    ...mapGetters(['userAlias', 'deadline', 'combinations']),
    isPrioritised () {
      return this.priorities.length === this.combinations.length
    }
  },
  methods: {
    ...mapActions('student', ['setIsConfirmed', 'setSignature']),
    onConfirm () {
      const { userAlias } = this
      this.setIsConfirmed({ userAlias, isConfirmed: true })
      this.modal.hide()
    },
    onClickConfirmButton () {
      this.modal.show()
    }
  }
}
</script>
