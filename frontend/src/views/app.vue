<template>
  <div class="main">
    <el-tabs v-model="editableTabsValue" type="card" @tab-remove="removeTab">
      <el-tab-pane key="Home" label="Home" name="Home">
        <HomeTab @onOpenConn="onOpenConn" />
      </el-tab-pane>
      <el-tab-pane
        v-for="(item, key) in ConnsOpenedData"
        :key="key"
        :label="key"
        :name="key"
        closable
        lazy
      >
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script>
import HomeTab from "./home_tab.vue";
import DataQueryTab from "./data_query_tab.vue";
import { ConnsOpened } from "../../wailsjs/go/main/App";

export default {
  components: { HomeTab, DataQueryTab },
  data() {
    return {
      ConnsOpenedData: [],
    };
  },
  computed: {},
  methods: {
    log(val) {},
    async ConnsOpened(key) {
      this.ConnsOpenedData = await ConnsOpened();
      this.editableTabsValue = key;

      setTimeout(() => {
        refreshTabEvent();
      }, 500);
    },
    onOpenConn(key) {
      setTimeout(() => {
        this.ConnsOpened(key);
      }, 500);
    },

    removeTab(targetName) {
      delete this.ConnsOpenedData[targetName];
      this.editableTabsValue = Object.keys(this.ConnsOpenedData)[0];
      if (!this.editableTabsValue) {
        this.editableTabsValue = "Home";
      }
    },
  },

  created() {},
  mounted() {},
};
</script>
