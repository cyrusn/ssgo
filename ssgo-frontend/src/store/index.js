// root state handle event for auth and user information
import student from "@/store/modules/student";
import { createStore } from "vuex";
import students from "@/store/modules/students";
import subject from "@/store/modules/subject";

import { checkError, alertMessage } from "@/store/helpers";
import _ from "lodash";

// https://stackoverflow.com/questions/38552003/how-to-decode-jwt-token-in-javascript-without-using-a-library#answer-38552302
function parseJWT(jwt) {
  const base64Url = jwt.split(".")[1];
  const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
  const jsonPayload = decodeURIComponent(
    atob(base64)
      .split("")
      .map(function (c) {
        return "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2);
      })
      .join("")
  );

  return JSON.parse(jsonPayload);
}

export default createStore({
  modules: {
    student,
    students,
    subject,
  },
  state: {
    jwt: "",
    message: "",
    config: null,
  },
  getters: {
    _info: (state) => (state.jwt ? parseJWT(state.jwt) : {}),
    role: (_, getters) => getters._info.Role || "",
    name: (_, getters) => getters._info.Name || "",
    cname: (_, getters) => getters._info.Cname || "",
    userAlias: (_, getters) => getters._info.UserAlias || "",
    expireAt: (_, getters) => {
      const exp = getters._info.exp;
      if (!exp || typeof exp !== 'number') return new Date();
      return new Date(exp * 1000);
    },
    config: (state) => state.config || {},
    schoolName: (state) => state.config ? state.config.schoolName : "",
    systemTitle: (state) => state.config ? state.config.systemTitle : "選科系統",
    schoolWebsite: (state) => state.config ? state.config.schoolWebsite : "",
    committeeInCharge: (state) => state.config ? state.config.committeeInCharge : "",
    committeeWebsite: (state) => state.config ? state.config.committeeWebsite : "",
    schoolYear: (state) => state.config ? state.config.schoolYear : "",
    deadline: (state) => state.config ? state.config.deadline : "",
    introductionMarkdown: (state) => state.config ? state.config.introMarkdown : "",
    instructionMarkdown: (state) => state.config ? state.config.instrMarkdown : "",
    notAcceptMarkdown: (state) => state.config ? state.config.notAcceptMarkdown : "",
    confirmMarkdown: (state) => state.config ? state.config.confirmMarkdown : "",
    isMock: (state) => state.config ? state.config.isMock : false,
    isAcceptResponses: (state) => state.config ? state.config.isAcceptResponses : true,
    warningTime: (state) => state.config ? state.config.warningTime : 50,
    subjects: (state) => {
      if (!state.config || !state.config.subjectsJSON) return [];
      try {
        return JSON.parse(state.config.subjectsJSON);
      } catch { return []; }
    },
    combinations: (state) => {
      if (!state.config || !state.config.combinationsJSON) return [];
      try {
        return JSON.parse(state.config.combinationsJSON);
      } catch { return []; }
    },
  },
  mutations: {
    updateJWT(state, jwt) {
      state.jwt = jwt;
    },
    updateMessage(state, message) {
      state.message = message;
    },
    updateConfig(state, config) {
      state.config = config;
    },
  },
  actions: {
    fetchConfig({ commit }) {
      return fetch("./api/config")
        .then((res) => {
          if (res.status === 503) {
            throw new Error("no_active_cohort");
          }
          return checkError(res);
        })
        .then((res) => res.json())
        .then((json) => {
          commit("updateConfig", json);
          return json;
        });
    },
    login({ commit, dispatch, getters }, { userAlias, password }) {
      fetch("./api/auth/login", {
        method: "POST",
        body: JSON.stringify({ userAlias, password }),
      })
        .then(checkError)
        .then((res) => res.text())
        .then((text) => {
          return text;
        })
        .then((text) => commit("updateJWT", text))
        .then(() => {
          const { role } = getters;
          switch (true) {
            case role === "STUDENT":
              dispatch("student/get");
              break;
            case _.includes(["TEACHER", "ADMIN"], role):
              dispatch("students/get");
              break;
          }
        })
        .catch(alertMessage);
    },
    refreshJWT({ state, commit }) {
      fetch("./api/auth/refresh/jwt", {
        method: "GET",
        headers: { jwt: state.jwt },
      })
        .then(checkError)
        .then((res) => res.text())
        .then((text) => commit("updateJWT", text))
        .catch(alertMessage);
    },
  },
});
