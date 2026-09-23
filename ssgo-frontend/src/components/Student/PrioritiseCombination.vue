<template>
  <div class="card" v-if="!isConfirmed">
    <h5 class="card-header bg-light">尚未編排的選科組合</h5>
    <div class="card-body">
      <draggable
        v-model="localAvailable"
        item-key="id"
        ghost-class="ghost"
        class="row"
        group="subject"
        :delay="100"
        :delayOnTouchOnly="true"
      >
        <template #item="{element, index}">
          <ul class="list-group list-group-horizontal col-md-4 col-6 mb-1 px-1">
            <subject-group
              :element="element"
              :index="index"
              :movable="!isConfirmed"
              name="available"
            />
          </ul>
        </template>
      </draggable>
    </div>
  </div>
  <div class="card my-4">
    <h5 class="card-header">
      已{{ isConfirmed ? '確定' : '編排' }}的選科組合次序
    </h5>
    <div class="card-body">
      <draggable
        v-model="localPrioritised"
        item-key="id"
        class="row"
        group="subject"
        :delay="100"
        :delayOnTouchOnly="true"
        @start="dragging = !isConfirmed"
      >
        <template #item="{element, index}">
          <div
            class="list-group list-group-horizontal col-md-4 col-6 mb-1 px-1"
          >
            <subject-group
              :element="element"
              :index="index"
              :movable="!isConfirmed"
              :isConfirmed="isConfirmed"
              name="prioritised"
            />
          </div>
        </template>
      </draggable>
    </div>
  </div>
</template>

<script>
import SubjectGroup from '@/components/Student/SubjectGroup'
import Draggable from 'vuedraggable'

import _ from 'lodash'
import { mapState, mapActions, mapGetters } from 'vuex'

export default {
  components: {
    SubjectGroup,
    Draggable
  },
  computed: {
    ...mapState('student', ['priorities', 'isConfirmed', 'timestamp']),
    ...mapGetters(['combinations']),
    localAvailable: {
      get () {
        const { priorities, combinations } = this
        return _.filter(combinations, comb => !_.includes(priorities, comb.id))
      },
      set () {
        // empty setter as localPrioritised handles the update
      }
    },
    localPrioritised: {
      get () {
        const { priorities, combinations } = this
        return _.map(priorities, id => _.find(combinations, { id })).filter(Boolean)
      },
      set (value) {
        const ids = value.map(e => e.id)
        this.updatePriorities(_.uniq(ids))
      }
    }
  },
  methods: {
    ...mapActions('student', ['updatePriorities'])
  }
}
</script>

<style>
.sortable-ghost,
.ghost {
  opacity: 0.2;
}
</style>
