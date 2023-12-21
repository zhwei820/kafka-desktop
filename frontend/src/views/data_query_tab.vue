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
      <multipane style="width: 100%" class="custom-resizer">
        <el-col
          :style="{ minWidth: '200px', width: '350px', maxWidth: '500px' }"
        >
          <div style="height: 90vh; overflow: scroll">
            <el-tooltip content="ADD DB" placement="top">
              <el-button
                style="margin-left: 20px; font-size: 20px"
                type="text"
                @click="() => addDB()"
                icon="el-icon-plus"
              >
              </el-button>
            </el-tooltip>
            <el-tooltip content="ADD TABLE" placement="top">
              <el-button
                style="margin-left: 20px; font-size: 20px"
                type="text"
                @click="() => addTable()"
                icon="el-icon-document-add"
              >
              </el-button>
            </el-tooltip>
            <el-tooltip content="REFRESH" placement="top">
              <el-button
                style="margin-left: 20px; font-size: 20px; margin-left: 20px"
                type="text"
                @click="() => refresh()"
                icon="el-icon-refresh"
              >
              </el-button>
            </el-tooltip>
            <el-tooltip content="DIFF" placement="top">
              <el-button
                style="margin-left: 20px; font-size: 20px; margin-left: 20px"
                type="text"
                @click="() => startTableDiff()"
                icon="el-icon-s-help"
              >
              </el-button>
            </el-tooltip>

            <el-tooltip content="CloseAll" placement="top">
              <el-button
                style="margin-left: 20px; font-size: 20px; margin-left: 20px"
                type="text"
                @click="() => startTableCloseAll()"
                icon="el-icon-close"
              >
              </el-button>
            </el-tooltip>

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
                  <!-- <el-tooltip :content="'DELETE TABLE'" placement="top">
                  <el-button
                    type="text"
                    v-if="!data.children"
                    @click="() => startTableDelete(data)"
                    icon="el-icon-delete"
                  >
                  </el-button>
                </el-tooltip> -->
                  <!-- <el-tooltip
                  :content="'TRUNCATE ' + (data.children ? 'DB' : 'TABLE')"
                  placement="top"
                >
                  <el-button
                    type="text"
                    @click="() => startTableTruncate(data)"
                    :icon="
                      data.children
                        ? 'el-icon-folder-delete'
                        : 'el-icon-document-delete'
                    "
                  >
                  </el-button>
                </el-tooltip> -->
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
        <multipane-resizer></multipane-resizer>
        <div :style="{ flexGrow: 1, maxWidth: '80vw' }">
          <el-col :span="24">
            <el-tabs
              class="fixed-tabs"
              v-model="editableTabsValue"
              type="card"
              @tab-remove="removeTab"
            >
              <el-tab-pane
                v-for="(item, index) in ConnsOpenedData[connName].OpenedConf
                  .SQLs"
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
                    <div
                      style="
                        display: flex;
                        justify-content: flex-end;
                        margin-right: 20px;
                      "
                    >
                      <el-button
                        size="mini"
                        plain
                        type="primary"
                        @click="
                          showCreateTable(
                            '`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][3] +
                              '`.`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][0] +
                              '`'
                          )
                        "
                      >
                        SHOW CREATE TABLE
                      </el-button>
                      <el-button
                        size="mini"
                        plain
                        type="primary"
                        @click="
                          exportCSV(
                            tableData[index],
                            '`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][3] +
                              '`.`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][0] +
                              '`'
                          )
                        "
                      >
                        EXPORT CSV
                      </el-button>
                      <el-button
                        size="mini"
                        plain
                        type="primary"
                        @click="
                          exportSQL(
                            tableData[index],
                            '`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][3] +
                              '`.`' +
                              ConnsOpenedData[connName].OpenedConf.SQLs[
                                index
                              ][0] +
                              '`'
                          )
                        "
                      >
                        EXPORT SQL
                      </el-button>
                    </div>

                    <multipane
                      style="width: 100%; height: 85vh"
                      layout="horizontal"
                      v-on:paneResizeStop="paneresizestop"
                      v-on:paneResize="paneresizestop"
                    >
                      <div>
                        <Editor
                          :keywords="
                            hintKeys[connName]
                              ? hintKeys[connName][index]
                                ? hintKeys[connName][index]
                                : []
                              : []
                          "
                          v-model="
                            ConnsOpenedData[connName].OpenedConf.SQLs[index][1]
                          "
                          @onSelectionChange="selectionChange"
                          @onUpdateValue="onUpdateValue"
                        />
                      </div>
                      <multipane-resizer></multipane-resizer>
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
                    </multipane>
                  </div>
                  <div v-else>
                    <TableDiff
                      v-if="index === 'TABLE-DIFF'"
                      @startDiff2="startDiff2"
                      @startDiff1="startDiff1"
                      :db_name="diff_db_nameVal"
                      :table_name="diff_table_nameVal"
                      :connName="diff_connNameVal"
                      :db_name2="diff_db_name2Val"
                      :table_name2="diff_table_name2Val"
                      :connName2="diff_connName2Val"
                      :diff2="diff2Val"
                      :diff1="diff1Val"
                    ></TableDiff>
                    <el-tabs type="border-card" v-else>
                      <el-tab-pane label="Setting">
                        <TableSettingSpreadsheet
                          :can-edit="true"
                          :table-data-prop="tableData[index] ?? []"
                          @applyTableSettingDiffData="applyTableSettingDiffData"
                          :table_name="
                            '`' +
                            ConnsOpenedData[connName].OpenedConf.SQLs[
                              index
                            ][3] +
                            '`.`' +
                            ConnsOpenedData[connName].OpenedConf.SQLs[
                              index
                            ][0] +
                            '`'
                          "
                          style="margin-bottom: 10px"
                        />
                      </el-tab-pane>
                      <el-tab-pane label="Index">
                        <TableIndexSpreadsheet
                          :can-edit="true"
                          :table-data-prop="tableDataIndex[index] ?? []"
                          :table-data-index-keys-prop="
                            tableDataIndexKeys[index] ?? []
                          "
                          @applyTableIndexDiffData="applyTableSettingDiffData"
                          :table_name="
                            '`' +
                            ConnsOpenedData[connName].OpenedConf.SQLs[
                              index
                            ][3] +
                            '`.`' +
                            ConnsOpenedData[connName].OpenedConf.SQLs[
                              index
                            ][0] +
                            '`'
                          "
                          style="margin-bottom: 10px"
                        />
                      </el-tab-pane>
                    </el-tabs>
                  </div>

                  <div
                    style="
                      margin-right: 10px;
                      margin-bottom: 0px;
                      z-index: 1000;
                      float: right;
                    "
                  >
                    <el-popover placement="right" width="1200" trigger="hover">
                      <div style="overflow: scroll; max-height: 700px">
                        <el-table :data="tableDataLog" border height="450">
                          <el-table-column
                            prop="status"
                            label="Status"
                            width="80"
                          >
                          </el-table-column>
                          <el-table-column prop="time" label="Time" width="180">
                          </el-table-column>
                          <el-table-column prop="sql" label="Sql">
                          </el-table-column>
                          <el-table-column
                            prop="result"
                            label="Result"
                            width="380"
                          >
                          </el-table-column>
                        </el-table>
                      </div>
                      <el-button type="success" slot="reference">
                        Action Logs
                      </el-button>
                    </el-popover>

                    <el-popover placement="right" width="1200" trigger="hover">
                      <div style="overflow: scroll; max-height: 700px">
                        <el-table
                          :data="tableDataSqlCollection"
                          border
                          height="450"
                        >
                          <el-table-column prop="name" label="Name" width="180">
                          </el-table-column>
                          <el-table-column prop="sql" label="Sql">
                          </el-table-column>
                        </el-table>
                      </div>
                      <el-button type="success" slot="reference">
                        Sql Collections
                      </el-button>
                    </el-popover>
                  </div>
                </div>
              </el-tab-pane>
            </el-tabs>
          </el-col>
        </div>
      </multipane>
    </el-row>

    <el-dialog title="Confirm SQL" :visible.sync="dialogVisible" width="50%">
      <el-input
        type="textarea"
        :rows="10"
        placeholder="请输入内容"
        v-model="sqlsToExecJoin"
      >
      </el-input>

      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="confirmApply">Confirm</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { Multipane, MultipaneResizer } from "vue-multipane";

import I18nMixins from "../comp/mixins/i18n-mixins";
import ThemeSwitchMixins from "../comp/mixins/theme-switch-mixins.js";
import clickoutside from "../comp/directives/clickoutside.js";
import Spreadsheet from "../comp/spreadsheet.vue";
import TableSettingSpreadsheet from "../comp/table-setting-spreadsheet.vue";
import TableIndexSpreadsheet from "../comp/table-index-spreadsheet.vue";
import TableDiff from "../comp/table-diff.vue";
import Editor from "../comp/editor.vue";
import {
  SetConnsOpened,
  Query,
  ConnsOpened,
  Execute,
  SaveFile,
  ShowCreateTable,
} from "../../wailsjs/go/main/App";
import { EventsOn, ClipboardSetText } from "../../wailsjs/runtime";
import { diffSqls } from "../utils";
import VueSimpleContextMenu from "../comp/vue-simple-context-menu.vue";

export default {
  directives: {
    "click-outside": clickoutside,
  },
  components: {
    Spreadsheet,
    VueSimpleContextMenu,
    TableSettingSpreadsheet,
    TableIndexSpreadsheet,
    TableDiff,
    Editor,
    Multipane,
    MultipaneResizer,
  },
  props: {
    connName: {
      type: String,
    },

    SelectedConn: {
      type: String,
    },

    diff_db_name: {
      type: String,
    },
    diff_table_name: {
      type: String,
    },
    diff_connName: {
      type: String,
    },
    diff_db_name2: {
      type: String,
    },
    diff_table_name2: {
      type: String,
    },
    diff_connName2: {
      type: String,
    },
    diff1: {
      type: Number,
    },
    diff2: {
      type: Number,
    },
  },
  mixins: [I18nMixins, ThemeSwitchMixins],
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
        // {
        //   name: "COPY TABLES CREATE",
        //   slug: "copy_tables_create",
        // },
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
        // {
        //   name: "COPY TABLES CREATE",
        //   slug: "copy_tables_create",
        // },
      ],
      diff1Val: 0,
      diff2Val: 0,

      diff_db_nameVal: "",
      diff_table_nameVal: "",
      diff_connNameVal: "",
      diff_db_name2Val: "",
      diff_table_name2Val: "",
      diff_connName2Val: "",

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
      tableDataLog: [],
      tableDataSqlCollection: [
        {
          sql: "SELECT *,CONCAT('KILL ', id, ';') FROM information_schema.processlist WHERE db = 'fixedincome' and state like '%lock%';",
          name: "query lock",
        },
        {
          sql: "show full processlist;",
          name: "show full processlist",
        },
        {
          sql: "show FULL columns from cell.cells;",
          name: "show FULL columns",
        },
      ],
      defaultProps: {
        children: "children",
        label: "label",
      },

      tableData: {},
      tableDataIndex: {},
      tableDataIndexKeys: {},
      dbtableColumns: {},
      selectedCode: "",
    };
  },
  watch: {
    diff1() {
      this.diff1Val = this.diff1;
    },
    diff2() {
      this.diff2Val = this.diff2;
    },

    diff_db_name() {
      this.diff_db_nameVal = this.diff_db_name;
    },
    diff_table_name() {
      this.diff_table_nameVal = this.diff_table_name;
    },
    diff_connName() {
      this.diff_connNameVal = this.diff_connName;
    },
    diff_db_name2() {
      this.diff_db_name2Val = this.diff_db_name2;
    },
    diff_table_name2() {
      this.diff_table_name2Val = this.diff_table_name2;
    },
    diff_connName2() {
      this.diff_connName2Val = this.diff_connName2;
    },
    editableTabsValue: {
      handler: function () {
        this.syncHintKeys(this.editableTabsValue);
      },
      immediate: true,
    },
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
    startDiff2() {
      this.diff2Val = 1;
      // this.diff_db_nameVal = "";
      // this.diff_table_nameVal = "";
      // this.diff_connNameVal = "";

      this.$emit("startDiff2", {
        diff2: this.diff2Val,
        diff_db_name2: this.diff_db_name2Val,
        diff_table_name2: this.diff_table_name2Val,
        diff_connName2: this.diff_connName2Val,
      });
    },
    startDiff1() {
      this.diff1Val = 1;
      // this.diff_db_name2Val = "";
      // this.diff_table_name2Val = "";
      // this.diff_connName2Val = "";

      this.$emit("startDiff1", {
        diff1: this.diff1Val,
        diff_db_name: this.diff_db_nameVal,
        diff_table_name: this.diff_table_nameVal,
        diff_connName: this.diff_connNameVal,
      });
    },

    async addDB() {
      this.startAddDb = true;
      this.sqlsToExec = ["CREATE SCHEMA `` DEFAULT CHARACTER SET UTF8MB4"];
      this.sqlsToExecJoin = this.sqlsToExec.join("\n\n");
      this.dialogVisible = true;
    },
    async addTable() {
      this.startAddDb = true;

      this.sqlsToExec = [
        `CREATE TABLE \`${
          this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema
        }\`.\`\` (
  id bigint(20) unsigned NOT NULL AUTO_INCREMENT,

  PRIMARY KEY (id)
) DEFAULT CHARSET=utf8mb4 comment 'no qa';  `,
      ];
      this.sqlsToExecJoin = this.sqlsToExec.join("\n\n");
      this.dialogVisible = true;
    },
    async refresh() {
      await this.refreshTree();

      this.$message.success("refresh success!");
    },
    showCreateTable(tableName) {
      ShowCreateTable(this.connName, tableName);
      this.$message.success("copy to clipboard!");
    },
    exportCSV(data, tableName) {
      let outputSata = convertToCSV(data);
      // outputSata
      this.export(outputSata, tableName, "csv");
    },
    exportSQL(data, tableName) {
      let outputSata = convertToSQL(tableName, data);
      // outputSata
      this.export(outputSata, tableName, "sql");
    },
    export(content, tableName, ext) {
      SaveFile(ext, tableName, content);
    },
    removeTab(targetName) {
      let tabs = this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs;

      delete tabs[targetName];
      this.editableTabsValue = Object.keys(tabs)[0];

      this.$forceUpdate();

      SetConnsOpened(
        this.connName,
        this.ConnsOpenedData[this.connName]?.OpenedConf
      );
    },

    async ConnsOpened() {
      this.ConnsOpenedData = await ConnsOpened();
      if (!this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs) {
        this.$forceUpdate();
        return;
      }
      this.editableTabsValue = Object.keys(
        this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs
      )[0];

      this.$forceUpdate();
      this.refreshColumnKeys();
    },
    async refreshColumnKeys() {
      let that = this;
      for (var key in this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs) {
        let tableName =
          this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs[key][0];
        if (that.dbtableColumns[tableName] || !tableName) {
          continue;
        }

        let sql =
          "SELECT column_name FROM information_schema.columns WHERE table_name ='" +
          tableName +
          "'";

        let res = await Query(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
          sql
        );

        that.dbtableColumns[tableName] = {};
        res.forEach((val) => {
          that.dbtableColumns[tableName][
            val.COLUMN_NAME ? val.COLUMN_NAME : val.column_name
          ] = "";
        });
      }
    },
    async startDiff(a) {
      console.log("==>startDiff", a);

      if (a.children) {
        return;
      }
      if (this.diff1Val === 1) {
        this.diff1Val = -1;
        this.diff_db_nameVal = a.db;
        this.diff_table_nameVal = a.label;
        this.diff_connNameVal = this.connName;

        this.$emit("startDiff1", {
          diff1: this.diff1Val,
          diff_db_name: this.diff_db_nameVal,
          diff_table_name: this.diff_table_nameVal,
          diff_connName: this.diff_connNameVal,
        });
      } else {
        this.diff2Val = -1;
        this.diff_db_name2Val = a.db;
        this.diff_table_name2Val = a.label;
        this.diff_connName2Val = this.connName;

        this.$emit("startDiff2", {
          diff2: this.diff2Val,
          diff_db_name2: this.diff_db_name2Val,
          diff_table_name2: this.diff_table_name2Val,
          diff_connName2: this.diff_connName2Val,
        });
      }
    },
    async startQuery(a) {
      console.log("===>startQuery", a);

      if (a.children) {
        return;
      }
      this.refreshColumnKeys();

      let sql =
        "select a.* from " +
        "`" +
        a.db +
        "`" +
        ".`" +
        a.label +
        "` a" +
        " where 1  order by id desc  limit 10;";

      let item = [a.label, sql, "", a.db];
      let sqlName = a.label;
      if (this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName]) {
        sqlName = sqlName + "__" + generateRandomString(5);
      }
      this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName] = item;

      this.editableTabsValue = sqlName;
      this.syncHintKeys(sqlName);

      this.$forceUpdate();
      SetConnsOpened(
        this.connName,
        this.ConnsOpenedData[this.connName]?.OpenedConf
      );
      this.Query(sqlName);
    },
    async syncHintKeys(sqlName) {
      if (
        !sqlName ||
        (this.hintKeys[this.connName] &&
          this.hintKeys[this.connName][sqlName] &&
          this.hintKeys[this.connName][sqlName].length > 0) ||
        !this.ConnsOpenedData[this.connName]
      ) {
        return;
      }
      if (
        !this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
          this.editableTabsValue
        ][0]
      ) {
        return;
      }

      this.hintKeys[this.connName] = {};

      let aLabel = sqlName.split("__")[0];
      let res;
      try {
        res = await Query(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
            this.editableTabsValue
          ][3],
          "SHOW COLUMNS FROM " +
            this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
              this.editableTabsValue
            ][0]
        );
      } catch (e) {
        console.log(
          "===>e",
          e,
          "SHOW COLUMNS FROM " +
            this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
              this.editableTabsValue
            ][0]
        );
        return;
      }
      let items = res
        .map((val) => {
          return val["Field"];
        })
        .sort();

      items.push(aLabel);
      this.hintKeys[this.connName][sqlName] = items;
    },
    async startTableSetting(a) {
      if (a.children) {
        return;
      }

      let sql = "";

      let item = [a.label, sql, "true", a.db];
      let sqlName = a.label + "_table_setting";
      if (this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName]) {
        sqlName = sqlName + "__" + generateRandomString(5);
      }
      this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName] = item;

      this.editableTabsValue = sqlName;

      this.reloadTableSetting();
      this.reloadTableIndex();
      // SetConnsOpened(
      //   this.connName,
      //   this.ConnsOpenedData[this.connName]?.OpenedConf
      // );
    },

    async startTableCloseAll(a) {
      this.$confirm("关闭全部?", "tip", {
        confirmButtonText: "Confirm",
        cancelButtonText: "Cancel",
        type: "warning",
      }).then(() => {
        this.ConnsOpenedData[this.connName].OpenedConf.SQLs = {};
      });
    },
    async startTableDiff(a) {
      let sql = "";
      let item = ["", sql, "true", ""];
      let sqlName = "TABLE-DIFF";
      if (this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName]) {
        return;
      }
      this.ConnsOpenedData[this.connName].OpenedConf.SQLs[sqlName] = item;

      this.diff1Val = 0;
      this.diff2Val = 0;
      this.diff_table_name = "";
      this.diff_table_name2 = "";

      this.editableTabsValue = sqlName;
    },
    async reloadTableSetting() {
      if (!this.connName) {
        return;
      }
      let res = await Query(
        this.connName,
        this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
          this.editableTabsValue
        ][3],
        "SHOW FULL COLUMNS FROM " +
          this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
            this.editableTabsValue
          ][0]
      );
      this.tableData[this.editableTabsValue] = res;
      // console.log("this.tableData[sqlName]", this.tableData[sqlName]);

      this.$forceUpdate();
    },

    async reloadTableIndex() {
      if (!this.connName) {
        return;
      }
      let res = await Query(
        this.connName,
        this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
          this.editableTabsValue
        ][3],
        "SHOW INDEX FROM " +
          this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
            this.editableTabsValue
          ][0]
      );
      this.tableDataIndex[this.editableTabsValue] = res;
      // console.log("this.tableData[sqlName]", this.tableData[sqlName]);

      try {
        res = await Query(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
            this.editableTabsValue
          ][3],
          "SHOW COLUMNS FROM " +
            this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
              this.editableTabsValue
            ][0]
        );
      } catch (e) {
        console.log(
          "===>e",
          e,
          "SHOW COLUMNS FROM " +
            this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
              this.editableTabsValue
            ][0]
        );

        return;
      }
      let items = res.map((val) => {
        return val["Field"];
      });
      this.tableDataIndexKeys[this.editableTabsValue] = items;

      this.$forceUpdate();
    },
    async startTableDelete(a) {
      if (a.children) {
        return;
      }

      let table =
        "`" +
        this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema +
        "`.`" +
        a.label +
        "`";
      let sql = "DROP TABLE  " + table;

      this.sqlsToExecJoin = [sql].join("\n\n");
      this.dialogVisible = true;
      this.startDelete = true;
    },
    async startDBDelete(a) {
      if (!a.children) {
        return;
      }

      let sql = "DROP DATABASE  " + a.label;

      this.sqlsToExecJoin = [sql].join("\n\n");
      this.dialogVisible = true;
      this.startDelete = true;
    },
    async toggleDbSTAR(a) {
      if (a.children) {
        if (
          this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema.includes(
            a.label
          )
        ) {
          const index = this.ConnsOpenedData[
            this.connName
          ].OpenedConf.StaredSchema.indexOf(a.label); // Find the index of the item
          if (index !== -1) {
            this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema.splice(
              index,
              1
            ); // Remove the item at the found index
          }
        } else {
          this.ConnsOpenedData[this.connName].OpenedConf.StaredSchema.unshift(
            a.label
          );
        }
        this.adjustTree();

        SetConnsOpened(
          this.connName,
          this.ConnsOpenedData[this.connName]?.OpenedConf
        );
      }
    },
    async startTableTruncate(a) {
      if (a.children) {
        this.sqlsToExecJoin = a.children
          .map((val) => {
            let table =
              "`" +
              this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema +
              "`.`" +
              val.label +
              "`";
            let sql = "TRUNCATE TABLE  " + table + ";";
            return sql;
          })
          .join("\n\n");
        this.dialogVisible = true;
        return;
      }

      let table =
        "`" +
        this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema +
        "`.`" +
        a.label +
        "`";
      let sql = "TRUNCATE TABLE  " + table + ";";

      this.sqlsToExecJoin = [sql].join("\n\n");
      this.dialogVisible = true;
    },
    handleNodeClick(a) {
      if (!a.children) {
        return;
      }
      this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema = a.label;

      SetConnsOpened(
        this.connName,
        this.ConnsOpenedData[this.connName]?.OpenedConf
      );
    },
    onQuery(sqlName) {
      this.Query(sqlName);
    },
    apply(a) {
      this.sqlsToExec = diffSqls(a.diffData, a.table_name);
      this.sqlsToExecJoin = this.sqlsToExec.join("\n\n");
      this.dialogVisible = true;
    },
    applyTableSettingDiffData(a) {
      this.startReloadTableSetting = true;
      this.sqlsToExec = a["diffData"];
      this.sqlsToExecJoin = this.sqlsToExec.join("\n\n");
      this.dialogVisible = true;
    },
    async confirmApply() {
      this.sqlsToExec = this.sqlsToExecJoin.split("\n\n");
      const loading = this.$loading({
        lock: true,
        text: "Loading",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)",
      });

      setTimeout(() => {
        loading.close();
      }, 10000);

      let res;
      let err;
      try {
        res = await Execute(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
          this.sqlsToExec
            .map((statement) => statement.trim())
            .filter((val) => {
              return val;
            })
        );
      } catch (e) {
        res = e;
        err = e;
        this.$message.error({
          showClose: true,
          message: "error: " + this.sqlsToExec.join() + " error:" + e,
        });
        loading.close();
      }

      this.tableDataLog.unshift({
        time: formatDateTime(new Date()),
        sql: this.sqlsToExec.join(),
        result: "affected: " + res + " rows",
        status: err ? "failed" : "success",
      });

      if (err) {
        return;
      }

      this.$message({
        message: "affected " + res + " rows; " + this.sqlsToExec.join(),
        showClose: true,
        type: "success",
      });

      loading.close();
      this.dialogVisible = false;

      if (this.startDelete) {
        this.refreshTree();
      }
      if (this.startReloadTableSetting) {
        this.reloadTableSetting();
      }
      if (this.startAddDb) {
        this.refreshTree();
      }
      this.startReloadTableSetting = false;
      this.startDelete = false;
      this.startAddDb = false;
    },
    async Query(index) {
      let code = this.selectedCode
        ? this.selectedCode
        : this.ConnsOpenedData[this.connName].OpenedConf.SQLs[index][1];

      if (!code) {
        return;
      }
      let res;
      let err;

      const loading = this.$loading({
        lock: true,
        text: "Loading",
        spinner: "el-icon-loading",
        background: "rgba(0, 0, 0, 0.7)",
      });

      let fn = Query;

      const tmp = code.trim().toLowerCase();
      let update = false;
      if (!(tmp.startsWith("select") || tmp.startsWith("show"))) {
        fn = Execute;
        update = true;

        code = code
          .split(";")
          .map((statement) => statement.trim())
          .filter((statement) => statement && !statement.startsWith("--"));
      }
      // console.log("===>tmp", tmp, update, fn);

      setTimeout(() => {
        loading.close();
      }, 10000);

      try {
        res = await fn(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
          code
        );
      } catch (e) {
        res = e;
        err = e;
        this.$message.error({
          showClose: true,
          message: "error: " + code + " error:" + e,
        });
        loading.close();
      }
      loading.close();

      this.tableDataLog.unshift({
        time: formatDateTime(new Date()),
        sql: code,
        result: !update
          ? JSON.stringify(res).slice(0, 200)
          : "affected: " + res + " rows",
        status: err ? "failed" : "success",
      });

      if (err) {
        return;
      }

      if (!(tmp.startsWith("select") || tmp.startsWith("show"))) {
        this.$message.success({
          message: "affected " + res + " rows; " + code,
          showClose: true,
        });
        return;
      }
      console.log("--->tmp", tmp, res);
      if (tmp.startsWith("select")) {
        this.canEdit = true;
        if (!res) {
          let tableName =
            this.ConnsOpenedData[this.connName]?.OpenedConf.SQLs[index][0];
          console.log(
            "111===>tableName",
            tableName,
            this.dbtableColumns,
            this.dbtableColumns[tableName]
          );

          if (this.dbtableColumns[tableName]) {
            res = [this.dbtableColumns[tableName]];
          }
        }
      } else {
        this.canEdit = false;
      }
      if (res) {
        this.tableData[index] = res;
      }
      this.$forceUpdate();
    },
    selectionChange(val) {
      this.selectedCode = val;
    },
    onUpdateValue(val) {
      this.ConnsOpenedData[this.connName].OpenedConf.SQLs[
        this.editableTabsValue
      ][1] = val;

      SetConnsOpened(
        this.connName,
        this.ConnsOpenedData[this.connName]?.OpenedConf
      );
    },

    async refreshTree() {
      let res = await Query(
        this.connName,
        this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
        "show databases"
      );
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
        res = await Query(
          this.connName,
          this.ConnsOpenedData[this.connName].OpenedConf.SelectedSchema,
          "SHOW TABLES IN `" + val.label + "`"
        );

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

  mounted() {
    let that = this;
    EventsOn("onQuery", function () {
      if (that.SelectedConn == that.connName) {
        that.onQuery(that.editableTabsValue);
      }
    });
    EventsOn("EmitQuery", function () {
      if (that.SelectedConn == that.connName) {
        that.Query(that.editableTabsValue);
      }
    });
    EventsOn("SaveQuery", function () {
      SetConnsOpened(
        that.connName,
        that.ConnsOpenedData[that.connName]?.OpenedConf
      );
    });
  },
};

function generateRandomString(length) {
  const charset = "abcdefghijklmnopqrstuvwxyz0123456789";
  let randomString = "";

  for (let i = 0; i < length; i++) {
    const randomIndex = Math.floor(Math.random() * charset.length);
    randomString += charset[randomIndex];
  }

  return randomString;
}

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

.custom-resizer > .multipane-resizer {
  margin: 0;
  left: 0;
  position: relative;

  &:before {
    display: block;
    content: "";
    width: 3px;
    height: 40px;
    position: absolute;
    top: 50%;
    left: 50%;
    margin-top: -20px;
    margin-left: -1.5px;
    border-left: 1px solid #ccc;
    border-right: 1px solid #ccc;
  }

  &:hover {
    &:before {
      border-color: #999;
    }
  }
}
</style>
