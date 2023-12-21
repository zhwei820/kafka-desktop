<template>
  <div class="main">
    <vue-simple-context-menu
      :element-id="'DBContextMenuID' + connName"
      :options="optionsDB"
      style="background-color: darkgray !important"
      :ref="'DBContextMenu' + connName"
      @option-clicked="optionClicked"
    />

    <vue-simple-context-menu
      :element-id="'TableContextMenuID' + connName"
      :options="optionsTable"
      style="background-color: darkgray !important"
      :ref="'TableContextMenu' + connName"
      @option-clicked="optionClicked"
    />
    <el-row
      :style="{
        display: 'flex !important',
      }"
    >
      <el-col :style="{ minWidth: '200px', width: '350px', maxWidth: '500px' }">
        <div style="height: 90vh; overflow: scroll">
          <el-tree
            style="background-color: lightgray"
            :data="data"
            :props="defaultProps"
            node-key="label"
            lazy
            @node-click="handleNodeClick"
            highlight-current
            :default-expanded-keys="[
              ConnsOpenedData[connName].OpenedConf.SelectedSchema,
            ]"
            accordion
          >
            <div
              @contextmenu.prevent="
                data.children
                  ? handleDBContextClick($event, data)
                  : handleTableContextClick($event, data)
              "
              class="custom-tree-node"
              slot-scope="{ node, data }"
              style="
                display: flex;
                justify-content: space-between;
                width: 86%;
                margin-right: 20px;
              "
            >
              <div style="max-width: 68%; overflow: hidden">
                {{ node.label }}
              </div>

              <div v-if="diff1Val === 1 || diff2Val === 1">
                <el-tooltip :content="'DIFF TABLE'" placement="top">
                  <el-button
                    type="text"
                    v-if="!data.children"
                    @click="() => startDiff(data)"
                    icon="el-icon-help"
                  >
                  </el-button>
                </el-tooltip>
              </div>
              <div v-else>
                <el-tooltip :content="'QUERY TABLE'" placement="top">
                  <el-button
                    type="text"
                    v-if="!data.children"
                    @click="() => startQuery(data)"
                    icon="el-icon-s-promotion"
                  >
                  </el-button>
                </el-tooltip>
                <el-tooltip :content="'SETTING TABLE'" placement="top">
                  <el-button
                    type="text"
                    v-if="!data.children"
                    @click="() => startTableSetting(data)"
                    icon="el-icon-setting"
                  >
                  </el-button>
                </el-tooltip>

                <el-tooltip :content="'STAR'" placement="top">
                  <el-button
                    type="text"
                    @click="() => toggleDbSTAR(data)"
                    v-if="data.children"
                    :icon="
                      ConnsOpenedData[
                        connName
                      ].OpenedConf.StaredSchema.includes(data.label)
                        ? 'el-icon-star-on'
                        : 'el-icon-star-off'
                    "
                  >
                  </el-button>
                </el-tooltip>
                &nbsp;
              </div>
            </div>
          </el-tree>
        </div>
      </el-col>

      <div :style="{ flexGrow: 1, maxWidth: '80vw' }">
        <el-col :span="24">
          <el-tabs
            class="fixed-tabs"
            v-model="editableTabsValue"
            type="card"
            @tab-remove="removeTab"
          >
            <el-tab-pane
              v-for="(item, index) in ConnsOpenedData[connName].OpenedConf.SQLs"
              :key="index"
              :label="index"
              :name="index"
              closable
              lazy
            >
              <div>
                <div
                  v-if="!ConnsOpenedData[connName].OpenedConf.SQLs[index][2]"
                  style="height: 85vh"
                >
                  <Spreadsheet
                    style="
                      align-content: space-between;
                      display: flex;
                      flex-direction: column;
                      flex-grow: 1;
                      max-width: 80vw;
                      margin-bottom: 10px;
                    "
                    :maxHeight="460 + 300 - size"
                    :table-data-prop="tableData[index] ?? []"
                    :dbtable-columns="
                      dbtableColumns[
                        ConnsOpenedData[connName].OpenedConf.SQLs[index][0]
                      ]
                    "
                    :can-edit="canEdit"
                    @apply="apply"
                    :table_name="
                      '`' +
                      ConnsOpenedData[connName].OpenedConf.SQLs[index][3] +
                      '`.`' +
                      ConnsOpenedData[connName].OpenedConf.SQLs[index][0] +
                      '`'
                    "
                  />
                </div>
              </div>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </div>
    </el-row>
  </div>
</template>

<script>
import clickoutside from "../comp/directives/clickoutside.js";
import { ConnsOpened } from "../../wailsjs/go/main/App";
import { EventsOn, ClipboardSetText } from "../../wailsjs/runtime";
import VueSimpleContextMenu from "../comp/vue-simple-context-menu.vue";

export default {
  directives: {
    "click-outside": clickoutside,
  },
  components: {
    VueSimpleContextMenu,
  },
  props: {
    connName: {
      type: String,
    },

    SelectedConn: {
      type: String,
    },
  },
  data() {
    return {
      size: 300,
      optionsDB: [
        {
          name: "COPY TABLES NAME",
          slug: "copy_tables_name",
        },
        {
          type: "divider",
        },
        {
          name: "TRUNCATE DB",
          slug: "truncate_db",
        },
        {
          name: "DELETE DB",
          slug: "delete_db",
        },
      ],
      optionsTable: [
        {
          name: "COPY TABLE NAME",
          slug: "copy_table_name",
        },
        {
          type: "divider",
        },
        {
          name: "DELETE TABLE",
          slug: "delete_table",
        },
        {
          name: "TRUNCATE TABLE",
          slug: "truncate_table",
        },
      ],

      canEdit: true,
      startDelete: false,
      startReloadTableSetting: false,
      startAddDb: false,

      dialogVisible: false,
      sqlsToExec: [],
      sqlsToExecJoin: "",

      leftSize: 5,
      editableTabsValue: "",

      ConnsOpenedData: {},
      hintKeys: {},

      data: [],

      defaultProps: {
        children: "children",
        label: "label",
      },

      tableData: {},
      dbtableColumns: {},
    };
  },

  computed: {},
  methods: {
    handleDBContextClick(event, item) {
      this.$refs["DBContextMenu" + this.connName].showMenu(event, item);
    },
    handleTableContextClick(event, item) {
      this.$refs["TableContextMenu" + this.connName].showMenu(event, item);
    },

    optionClicked(event) {
      if (
        event.option.slug === "copy_tables_name" ||
        event.option.slug === "copy_table_name"
      ) {
        this.onCopyTableName(event.item);
        this.$message.success("copy to clipboard!");
      } else if (
        event.option.slug === "truncate_db" ||
        event.option.slug === "truncate_table"
      ) {
        this.startTableTruncate(event.item);
      } else if (event.option.slug === "delete_table") {
        this.startTableDelete(event.item);
      } else if (event.option.slug === "delete_db") {
        this.startDBDelete(event.item);
      }
    },
    onCopyTableName(data) {
      let tables = [data];
      if (data.children) {
        tables = data.children;
      }
      let tablesName = tables.map((val) => {
        return "`" + val.label + "`";
      });
      ClipboardSetText(tablesName.join(","));
    },
    onCopyTableCreate(dbName) {},

    async ConnsOpened() {
      this.ConnsOpenedData = await ConnsOpened();
      if (!this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs) {
        this.$forceUpdate();
        return;
      }

      this.$forceUpdate();
    },

    handleNodeClick(a) {
      if (!a.children) {
        return;
      }
      this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema = a.label;
    },
    onQuery(sqlName) {
      this.Query(sqlName);
    },

    async refreshTree() {
      // let res = await Query(
      //   this.connName,
      //   this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
      //   "show databases"
      // );
      let res = [];
      this.data = res
        .map((val) => {
          return { label: val.Database, children: [] };
        })
        .filter((val) => {
          return (
            val.label !== "PERFORMANCE_SCHEMA" &&
            val.label !== "METRICS_SCHEMA" &&
            val.label !== "INFORMATION_SCHEMA" &&
            val.label !== "information_schema" &&
            val.label !== "mysql" &&
            val.label !== "performance_schema"
          );
        });
      // this.data
      // StaredSchema
      this.adjustTree();

      this.data.forEach(async (val, index) => {
        // res = await Query(
        //   this.connName,
        //   this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
        //   "SHOW TABLES IN `" + val.label + "`"
        // );
        res = [];
        if (!res) {
          return;
        }
        this.data[index].children = res.map((val1) => {
          return { label: val1["Tables_in_" + val.label], db: val.label };
        });
      });
    },
    adjustTree() {
      this.data.sort((a, b) => a.label - b.label);

      const movedToFront = [];
      if (!this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema) {
        this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema = [];
      }
      const remaining = this.data.filter((item) => {
        if (
          this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema.includes(
            item.label
          )
        ) {
          movedToFront.push(item);
          return false;
        }
        return true;
      });
      // this.data.splice(0, this.data.length, ...movedToFront.concat(remaining));
      this.data = movedToFront.concat(remaining);
      this.$forceUpdate();
    },
    paneresizestop(pane, container, size) {
      this.size = parseInt(size);
      console.log("pane, container, size", pane, container, size);
    },
  },

  async created() {
    await this.ConnsOpened();
    await this.refreshTree();

    this.$forceUpdate();
  },

  mounted() {},
};

function formatDateTime(date) {
  const year = date.getFullYear().toString().slice(-2);
  const month = ("0" + (date.getMonth() + 1)).slice(-2);
  const day = ("0" + date.getDate()).slice(-2);
  const hours = ("0" + date.getHours()).slice(-2);
  const minutes = ("0" + date.getMinutes()).slice(-2);
  const seconds = ("0" + date.getSeconds()).slice(-2);

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}

function convertToCSV(arr) {
  const header = Object.keys(arr[0])
    .filter((val) => {
      return val !== "rowKey";
    })
    .join(","); // Extract the object keys as the CSV header
  const csv = arr.map((obj) => {
    delete obj["rowKey"];
    return Object.values(obj).join(",");
  }); // Convert objects to CSV rows
  return [header, ...csv].join("\n"); // Join the header and rows with newline characters
}

function convertToSQL(tableName, data) {
  if (!Array.isArray(data) || data.length === 0) {
    return "";
  }

  let columns = Object.keys(data[0]).filter((val) => {
    return val !== "rowKey";
  });

  const values = data.map((item) => {
    const rowValues = columns.map((col) => {
      const value = item[col];
      // Ensure values are properly escaped and formatted for SQL
      return typeof value === "string" ? `'${value}'` : value;
    });
    return `(${rowValues.join(", ")})`;
  });

  columns = columns.map((val) => "`" + val + "`");
  const insertSQL = `INSERT INTO ${tableName} (${columns.join(
    ", "
  )}) VALUES\n${values.join(",\n")};`;

  return insertSQL;
}
</script>

<style>
.resizable-col {
  border-left: 2px solid #409eff; /* Add a border to indicate resizing */
  resize: horizontal; /* Enable horizontal resizing */
  overflow: hidden;
  padding: 0;
  cursor: col-resize;
}
.is-focusable {
  background-color: rgb(234, 237, 235);
}
.is-expanded {
  background-color: white;
}
.is-expanded > div {
  background-color: white;
}
.is-current {
  background-color: #f0f7ff !important;
}
.el-button + .el-button {
  margin-left: 5px !important;
}
.el-tree-node__children > .el-tree-node {
  padding-top: 5px !important;
  padding-bottom: 5px !important;
}

.fixed-tabs {
  position: unset;
}

.custom-tree-node {
  align-items: center;
}
.el-tabs__header {
  margin-bottom: 0 !important;
  margin-top: 0 !important;
}
.el-tabs__content {
  padding-bottom: 0 !important;
}
</style>
