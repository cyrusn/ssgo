// modules students only handle event about list of students

import { checkError, alertMessage } from "@/store/helpers";
import _ from "lodash";

export default {
  namespaced: true,
  state: {
    students: [
      /*
        'userAlias': String
        'name': String,
        'cname': String,
        'classCode': String,
        'classNo': Number,
        'priorities': [Number],
        'is3X': Boolean,
        'isConfirmed': Boolean,
        'timestamp': Time
        */
    ],
    signatures: [
      /*
        'userAlias': String
        'isSigned': Boolean,
        'address': String,
       */
    ],
  },
  mutations: {
    update(state, students) {
      state.students = students;
    },
    updateSignatures(state, signatures) {
      state.signatures = signatures;
    },
  },
  actions: {
    get({ commit, rootState }) {
      fetch("./api/students", {
        method: "GET",
        headers: { jwt: rootState.jwt },
      })
        .then(checkError)
        .then((res) => res.json())
        .then((ss) => {
          // ss stands for students
          const list = _(ss)
            // .map(s => {
            //   const { userAlias } = s
            //   return Object.assign({}, _.find(students, { userAlias }), s)
            // })
            .orderBy("classNo")
            .value();
          commit("update", list);
        })
        .catch(alertMessage);

      fetch("./api/signatures", {
        method: "GET",
        headers: { jwt: rootState.jwt },
      })
        .then(checkError)
        .then((res) => res.json())
        .then((list) => {
          commit("updateSignatures", list);
        });
    },
    updateRanks({ rootState, dispatch }, json) {
      fetch("./api/students/rank", {
        method: "PUT",
        headers: { jwt: rootState.jwt },
        body: JSON.stringify(json),
      })
        .then(checkError)
        .then(() => {
          dispatch("get");
        })
        .catch(alertMessage);
    },
  },
};
