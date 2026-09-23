<template>
  <div class="card mb-2" ref="container">
    <div class="card-header"><h5>家長簽署</h5></div>
    <div class="card-body">
      <signature-image v-if="isConfirmed" :signature="signature" />
      <div v-else>
        <vue-signature-pad
          class="border border-1 bg-light card-img-top mx-auto
        my-3"
          :width="width"
          :height="height"
          ref="signaturePad"
          :options="options"
        />
        <div class="d-flex flex-wrap gap-2">
          <button class="btn text-light btn-primary" @click="clear">
            重新簽署
          </button>
          <button 
            v-if="hasDrawn && !isConfirmed && isPrioritized" 
            class="btn btn-danger" 
            @click="submitAll"
            :disabled="isSubmitting"
          >
            <font-awesome-icon v-if="isSubmitting" icon="sync-alt" spin />
            確認簽署及選科次序
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import VueSignaturePad from "vue3-signature-pad";
import SignatureImage from "@/components/Shared/SignatureImage";
import { mapState, mapActions, mapMutations, mapGetters } from "vuex";

export default {
  components: {
    VueSignaturePad,
    SignatureImage,
  },
  data() {
    return {
      width: 600,
      height: 300,
      hasDrawn: false,
      isSubmitting: false,
      options: {
        onEnd: () => {
          this.hasDrawn = true;
        },
      },
    };
  },
  mounted() {
    const { handleResize } = this;
    window.addEventListener("resize", handleResize);
    handleResize();

    this.$nextTick(() => {
      const { isConfirmed, signature, $refs } = this;
      if ($refs.signaturePad) {
        if (isConfirmed) {
          $refs.signaturePad.lockSignaturePad();
        }
        if (signature) {
          $refs.signaturePad.fromDataURL(signature);
        }
      }
    });
  },
  unmounted() {
    window.removeEventListener("resize", this.handleResize);
  },
  watch: {
    signature() {
      const { signature, $refs } = this;
      if (signature && $refs.signaturePad) {
        $refs.signaturePad.fromDataURL(signature);
      }
    },
  },
  computed: {
    ...mapState("student", [
      "priorities",
      "signature",
      "isSigned",
      "isConfirmed",
    ]),
    ...mapGetters(['confirmMarkdown', 'combinations']),
    isPrioritized() {
      return this.priorities.length === this.combinations.length;
    },
  },
  methods: {
    ...mapMutations("student", ["updateIsSigned"]),
    ...mapActions("student", ["setSignature", "setIsConfirmed"]),
    handleResize() {
      const { $refs, signature } = this;
      const width = $refs.container ? $refs.container.clientWidth : 600;
      this.width = Math.min(width * 0.9, 600);
      this.height = this.width * 0.5;
      if (this.$refs.signaturePad) {
        this.$refs.signaturePad.resizeCanvas();
        if (signature) {
          this.$refs.signaturePad.fromDataURL(signature);
        }
      }
    },
    async submitAll() {
      const { signaturePad } = this.$refs;
      if (!signaturePad) return;
      const { isEmpty, data } = signaturePad.saveSignature();
      
      if (isEmpty) {
        alert("請先簽署");
        return;
      }

      const msg = this.confirmMarkdown || "是否確定簽署及選科次序？一經確定後，資料一概不得更改。";
      if (!confirm(msg)) return;

      this.isSubmitting = true;
      try {
        await this.setSignature({ isSigned: true, data });
        const userAlias = this.$store.getters.userAlias;
        await this.setIsConfirmed({ userAlias, isConfirmed: true });
        this.hasDrawn = false;
      } catch (err) {
        console.error("Submission failed:", err);
      } finally {
        this.isSubmitting = false;
      }
    },
    clear() {
      const { signaturePad } = this.$refs;
      if (signaturePad) {
        signaturePad.clearSignature();
      }
      this.hasDrawn = false;
    },
  },
};
</script>
