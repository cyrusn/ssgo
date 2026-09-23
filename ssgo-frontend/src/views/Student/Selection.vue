<template>
  <div style="max-width: 720px" class="mx-auto">
    <introduction />
    <instruction />
    <highlight-button-group />
    <prioritise-combination />
    <signature-pad v-if="isPrioritised" />
    <warning-message />
  </div>
</template>

<script>
import PrioritiseCombination from '@/components/Student/PrioritiseCombination'
import SignaturePad from '@/components/Student/SignaturePad'

import Introduction from '@/components/Student/Introduction'
import Instruction from '@/components/Student/Instruction'
import HighlightButtonGroup from '@/components/Student/HighlightButtonGroup'
import WarningMessage from '@/components/Student/WarningMessage'

import { mapState, mapGetters } from 'vuex'

export default {
  components: {
    Introduction,
    Instruction,
    HighlightButtonGroup,
    SignaturePad,
    PrioritiseCombination,
    WarningMessage
  },
  computed: {
    ...mapState('student', ['priorities', 'isSigned']),
    ...mapGetters(['systemTitle', 'combinations']),
    isPrioritised() {
      return this.priorities.length === this.combinations.length
    }
  }
}
</script>
