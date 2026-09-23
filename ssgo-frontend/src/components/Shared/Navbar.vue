<template>
  <nav :class="navbarClass">
    <div class="container-fluid">
      <a :class="navbarBrandClass" :href="schoolWebsite"
        >{{ schoolName }}</a
      >
      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbar"
      >
        <span class="navbar-toggler-icon"></span>
      </button>

      <div class="show navbar-collapse" id="navbar" v-if="userAlias">
        <ul class="navbar-nav me-auto">
          <router-list-link v-for="n in Navs" :name="n.name" :key="n.name">
            <font-awesome-icon :icon="n.icon" /> {{ n.content }}
          </router-list-link>
          <li class="nav-item nav-link" @click="onLogout" style="cursor: pointer;">
            <font-awesome-icon icon="sign-out-alt" /> 登出
          </li>
        </ul>
        <span class="navbar-text me-2">
          <font-awesome-icon icon="user" />
          {{ display_cname }}
        </span>
        <button
          id="logoutCounter"
          class="btn btn-outline-danger"
          data-bs-toggle="tooltip"
          data-bs-placement="bottom"
          title="更新登出時間"
          @click="refreshJWT"
          @mouseover="toggle"
        >
          <font-awesome-icon icon="hourglass" />{{ " " }} 尚餘{{
            logoutFromNow
          }}
        </button>
      </div>
      <logout-alert />
    </div>
  </nav>
</template>

<script>
import _ from "lodash";

import RouterListLink from "@/components/Shared/RouterListLink.vue";
import LogoutAlert from "@/components/Shared/LogoutAlert";
import moment from "moment";

import { mapGetters, mapActions } from "vuex";
import { Tooltip } from "bootstrap";

const navbarClass = [
  "navbar",
  "navbar-expand-lg",
  "navbar-light",
  "bg-light",
  "px-lg-5",
  "text-secondary",
  "d-print-none",
];
const navbarBrandClass = ["navbar-brand", "text-secondary"];

export default {
  mounted() {
    const { updateLogoutFromNow } = this;
    this.timer = setInterval(() => {
      updateLogoutFromNow();
    }, 500);
  },
  beforeUnmount() {
    if (this.timer) clearInterval(this.timer);
  },
  data() {
    return {
      navbarClass,
      navbarBrandClass,
      logoutFromNow: "...",
      timer: null
    };
  },
  components: {
    RouterListLink,
    LogoutAlert,
  },
  computed: {
    ...mapGetters(["cname", "name", "role", "userAlias", "expireAt", "schoolName", "schoolWebsite"]),
    logoutCounter() {
      const el = document.getElementById("logoutCounter");
      if (!el) return null;
      return new Tooltip(el);
    },
    display_cname() {
      const { cname, role } = this;
      if (role === "STUDENT") {
        return cname + "同學";
      }
      if (role === "ADMIN") {
        return cname + "管理員";
      }
      if (role === "SUPERADMIN") {
        return "系統管理員";
      }
      return cname + "老師";
    },
    Navs() {
      const { role } = this;
      const navs = [
        {
          roles: ["STUDENT"],
          name: "selection",
          icon: "heart",
          content: "選科",
        },
        {
          roles: ["TEACHER", "ADMIN"],
          name: "list",
          icon: "list-alt",
          content: "列表",
        },
        {
          roles: ["ADMIN"],
          name: "rank",
          icon: "trophy",
          content: "名次",
        },
        {
          roles: ["ADMIN"],
          name: "allocation",
          icon: "sitemap",
          content: "分科",
        },
        {
          roles: ["ADMIN"],
          name: "statistics",
          icon: "chart-bar",
          content: "統計",
        },
      ];
      return navs.filter((n) => _.includes(n.roles, role));
    },
  },
  methods: {
    ...mapActions(["refreshJWT"]),
    onLogout() {
      window.location.reload(true);
    },
    toggle() {
      if (this.logoutCounter) {
        this.logoutCounter.show();
      }
    },
    updateLogoutFromNow() {
      if (!this.expireAt || isNaN(this.expireAt.getTime())) {
        this.logoutFromNow = "...";
        return;
      }
      moment.updateLocale("zh-hk", {
        relativeTime: {
          mm: "%d分鐘",
        },
      });
      this.logoutFromNow = moment(this.expireAt).fromNow(true);
    },
  },
};
</script>
