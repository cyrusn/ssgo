<template>
  <div>
    <div class="alert alert-info">
      <b>已獲分發選修科人數 / 總人數</b>
      <span class="badge bg-danger rounded-pill ms-2"
        >{{ allocationResults.length }} / {{ filteredStudents.length }}</span
      >
    </div>
    <div class="alert alert-warning">
      <h6>已滿額選修科目</h6>
      <span
        v-for="s in occupiedSubjects"
        :key="s"
        class="badge bg-warning text-dark me-1"
        >{{ findSubject(s).code.toUpperCase() }}</span
      >
    </div>
    <table class="table table-border text-center table-striped">
      <tr>
        <th>獲派志願</th>
        <th>人次</th>
      </tr>

      <tr v-for="(value, key) in statistic" :key="key" class="border">
        <td class="border">{{ key }}</td>
        <td class="border">{{ value }}</td>
      </tr>
    </table>

    <subject-occupied-counters :counters="counters" />
    <download-results-button-group :results="allocationResults" />
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";
import _ from "lodash";
import SubjectOccupiedCounters from "@/components/Admin/SubjectOccupiedCounters";
import DownloadResultsButtonGroup from "@/components/Admin/DownloadResultsButtonGroup";
import { removeCombinationWithX3Subject } from "@/components/Shared/helpers";

export default {
  components: {
    SubjectOccupiedCounters,
    DownloadResultsButtonGroup,
  },
  data() {
    return {
      counters: {},
      occupiedSubjects: [],
      statistic: {
        "1": 0,
        "<=2": 0,
        "<=3": 0,
        "<=10": 0,
        "<=20": 0,
        ">20": 0,
      },
    };
  },
  computed: {
    ...mapState("students", ["students"]),
    ...mapState("subject", ["capacities"]),
    ...mapGetters(["combinations", "subjects"]),
    defaultCounters() {
      const counters = {};
      this.subjects.forEach(s => { counters[s.code] = 0 });
      return counters;
    },
    filteredStudents() {
      const { filter, students } = this;
      return filter(students);
    },
    allocationResults() {
      const {
        resetCounters,
        allocate,
        filteredStudents,
        resetOccupiedSubjects,
        updatePreferenceStatistic,
      } = this;
      resetOccupiedSubjects();
      resetCounters();
      const allocatedStudents = allocate(filteredStudents);
      const statistic = {
        "1": 0,
        "<=2": 0,
        "<=3": 0,
        "<=10": 0,
        "<=20": 0,
        ">20": 0,
      };
      allocatedStudents.forEach((s) => {
        const pref = s.preference || 0;
        Object.keys(statistic).forEach((key) => {
          switch (key) {
            case "1":
              if (pref == 1) statistic[key] += 1;
              break;
            case "<=2":
              if (pref > 0 && pref <= 2) statistic[key] += 1;
              break;
            case "<=3":
              if (pref > 0 && pref <= 3) statistic[key] += 1;
              break;
            case "<=10":
              if (pref > 0 && pref <= 10) statistic[key] += 1;
              break;
            case "<=20":
              if (pref > 0 && pref <= 20) statistic[key] += 1;
              break;
            case ">20":
              if (pref > 20) statistic[key] += 1;
              break;
            default:
              break;
          }
        });
      });
      updatePreferenceStatistic(statistic);
      return _.filter(allocatedStudents, (s) => s.preference > 0);
    },
  },
  methods: {
    updatePreferenceStatistic(statistic) {
      this.statistic = Object.assign({}, statistic);
    },
    findSubject(code) {
      return _.find(this.subjects, { code }) || { code: code.toUpperCase() };
    },
    resetOccupiedSubjects() {
      this.occupiedSubjects = [];
    },
    pushOccupiedSubjects(id) {
      const { occupiedSubjects } = this;
      if (_.includes(occupiedSubjects, id)) return;
      occupiedSubjects.push(id);
    },
    updateOccupiedSubjects() {
      const { capacities, pushOccupiedSubjects } = this;
      _.forOwn(this.counters, (v, k) => {
        if (v === capacities[k]) {
          pushOccupiedSubjects(k);
        }
      });
    },
    resetCounters() {
      this.counters = Object.assign({}, this.defaultCounters);
    },
    updateCounter(code) {
      this.counters[code] += 1;
    },
    getSubjectCodes(id) {
      const combo = _.find(this.combinations, { id });
      return combo ? combo.subjects : [];
    },
    isAvailable(id) {
      const { getSubjectCodes, counters, capacities } = this;
      const codes = getSubjectCodes(id);
      if (!codes || codes.length < 2) return false;
      const subject1 = codes[0];
      const subject2 = codes[1];
      const validSubject1 = counters[subject1] < (capacities[subject1] || 0);
      const validSubject2 = counters[subject2] < (capacities[subject2] || 0);

      return validSubject1 && validSubject2;
    },
    makeOffers(priorities) {
      let offers;
      let orders;
      let preference;
      const {
        isAvailable,
        getSubjectCodes,
        updateCounter,
        updateOccupiedSubjects,
      } = this;

      _.forEach(priorities, (id, i) => {
        if (isAvailable(id)) {
          const codes = getSubjectCodes(id);
          codes.forEach(updateCounter);

          offers = _.zipObject(["subject1", "subject2"], codes);
          orders = _.zipObject(
            ["subject1", "subject2"],
            codes.map((code) => this.counters[code]),
          );
          preference = i + 1;
          // break lodash forEach loop
          return false;
        }
      });
      updateOccupiedSubjects();
      return {
        offers: offers || {},
        preference: preference || 0,
        orders: orders || {},
      };
    },
    filter(students) {
      return _(students)
        .filter({ isConfirmed: true })
        .filter((s) => (s.priorities || []).length === this.combinations.length)
        .filter((s) => s.rank > 0)
        .orderBy("rank")
        .value();
    },
    allocate(students) {
      const { makeOffers, combinations } = this;
      return _(students)
        .map((s) => {
          if (!s.isX3) return s;
          const priorities = removeCombinationWithX3Subject(
            s.priorities,
            "hmsc",
            combinations
          );
          return Object.assign({}, s, { priorities });
        })
        .map((s) => Object.assign({}, s, makeOffers(s.priorities)))
        .orderBy(["classCode", "classNo"])
        .value();
    },
  },
};
</script>
