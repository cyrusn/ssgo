<template>
  <table class="table table-striped mt-2">
    <thead class="thead-light">
      <tr>
        <th>學生編號</th>
        <th>班別</th>
        <th>學號</th>
        <th>英文姓名</th>
        <th>中文姓名</th>
        <th>已確定 ／ 尚未確定</th>
        <th>修讀健管</th>
        <th>確定時間</th>
        <th style="max-width:10%">家長簽署</th>
        <th v-if="role === 'ADMIN'">級名次</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="s in list" :key="s.id">
        <td>{{ s.userAlias }}</td>
        <td>{{ s.classCode }}</td>
        <td>{{ s.classNo }}</td>
        <td>{{ s.name }}</td>
        <td>{{ s.cname }}</td>
        <td>
          <span
            v-if="role === 'TEACHER'"
            :class="!s.isConfirmed ? 'text-danger' : ''"
            >{{ s.isConfirmed ? "已確定" : "尚未確定" }}</span
          >
          <input
            v-else
            type="checkbox"
            :checked="s.isConfirmed"
            @click="
              setIsConfirmed({
                userAlias: s.userAlias,
                isConfirmed: !s.isConfirmed,
              })
            "
          />
          {{ s.isConfirmed ? "已確定" : "尚未確定" }}
        </td>
        <td>
          <span
            v-if="role === 'TEACHER'"
            :class="!s.isX3 ? 'text-danger' : ''"
            >{{ s.isX3 ? "修讀" : "不適用" }}</span
          >
          <input
            v-else
            type="checkbox"
            class="btn"
            :class="s.isX3 ? 'btn-success' : 'btn-danger'"
            @click="setIsX3({ userAlias: s.userAlias, isX3: !s.isX3 })"
            :checked="s.isX3"
          />
          {{ s.isX3 ? "修讀" : "不適用" }}
        </td>
        <td>
          <formatted-datetime
            v-if="s.timestamp.Valid"
            :datetime="s.timestamp.Time"
            format="lll"
          />
        </td>
        <td>
          <zoom-signature :src="getSignatureAddressByUserAlias(s.userAlias)" />
        </td>
        <td
          v-if="role === 'ADMIN'"
          :class="s.rank === 0 ? 'text-danger' : 'text-success'"
        >
          {{ s.rank === 0 ? "尚未提供" : s.rank }}
        </td>
      </tr>
    </tbody>
  </table>
</template>

<script>
import { mapGetters, mapActions } from "vuex";
import FormattedDatetime from "@/components/Teacher/FormattedDatetime";
import ZoomSignature from "@/components/Teacher/ZoomSignature";
import _ from "lodash";

export default {
  props: ["list", "signatures"],
  components: { FormattedDatetime, ZoomSignature },
  computed: {
    ...mapGetters(["role"]),
  },
  methods: {
    ...mapActions("student", ["setIsConfirmed", "setIsX3"]),
    getSignatureAddressByUserAlias(userAlias) {
      const { signatures } = this;
      const result = _.find(signatures, { userAlias });
      if (result == undefined) {
        return "";
      }

      if ("address" in result) {
        return result.address;
      }
      return "";
    },
  },
};
</script>

<style>
.zoom:hover {
  transform: scale(6);
}
</style>
