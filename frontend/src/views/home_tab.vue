<template>
  <div style="padding: 20px">
    <el-button type="primary" @click="dialogVisible = true">Add</el-button>

    <el-dialog title="Connection" :visible.sync="dialogVisible" width="50%">
      <el-form :model="c" label-position="right" label-width="100px">
        <el-form-item label="Name">
          <el-input
            v-model="c.Name"
            :disabled="c.Type !== ''"
            autocomplete="off"
          ></el-input>
        </el-form-item>
        <el-form-item label="Host">
          <el-input v-model="c.Host" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Port">
          <el-input v-model="c.Port" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="User">
          <el-input v-model="c.User" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Password">
          <el-input v-model="c.Password" autocomplete="off"></el-input>
        </el-form-item>
        <el-form-item label="Database">
          <el-input v-model="c.Database" autocomplete="off"></el-input>
        </el-form-item>
      </el-form>

      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">Cancel</el-button>
        <el-button type="primary" @click="addTab(c.Name, c)">Confirm</el-button>
      </span>
    </el-dialog>
    <div class="grid-container">
      <el-card class="box-card" v-for="(val, key) in DBConf" :key="key">
        <div slot="header">
          <div style="display: flex; justify-content: space-between">
            <h2>{{ DBConf[key].Name }}</h2>

            <el-button type="primary" plain @click="OpenConn(key)">
              OPEN
            </el-button>
          </div>
        </div>

        <div
          class="text item"
          style="
            display: flex;
            flex-direction: column;
            justify-content: space-between;
            height: 120px;
          "
        >
          <div>
            <strong>Host:</strong>
            {{ " " + DBConf[key].Host }}
          </div>
          <div style="display: flex; justify-content: flex-end">
            <el-button small type="text" @click="startEdit(DBConf[key])"
              >edit</el-button
            >
            <el-button small type="text" @click="startDuplicate(DBConf[key])"
              >duplicate</el-button
            >
            <el-button
              size="small"
              plain
              type="danger"
              @click="deleteDBConf(key)"
              >delete</el-button
            >
          </div>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script>
import { DBConf } from "../../wailsjs/go/main/App";
import { OpenConn } from "../../wailsjs/go/main/App";
import { SetDBConf } from "../../wailsjs/go/main/App";
import { ClipboardSetText } from "../../wailsjs/runtime/runtime";
export default {
  components: {},
  data() {
    return {
      dialogVisible: false,

      DBConf: {},
      c: {
        Type: "",
        Name: "",
        Host: "",
        Port: "",
        User: "",
        Password: "",
        Database: "",
      },
    };
  },
  computed: {},
  methods: {
    startEdit(val) {
      this.c = val;
      this.c.Type = "kafka";
      this.dialogVisible = true;
    },
    startDuplicate(val) {
      this.c = { ...val };
      this.c.Type = "";
      this.dialogVisible = true;
    },
    deleteDBConf(val) {
      this.$confirm("Delete Instance Config?", "Warning", {
        confirmButtonText: "Confirm",
        cancelButtonText: "Cancel",
        type: "warning",
      }).then(async () => {
        // delete this.DBConf[val];
        this.$delete(this.DBConf, val);

        await SetDBConf(this.DBConf);

        this.$message({
          type: "success",
          message: "Delete Success!",
        });

        this.$forceUpdate();
      });
    },
    removeTab(targetName) {
      let tabs = this.DBConf;

      delete tabs[targetName];
      this.editableTabsValue = Object.keys(tabs)[0];

      this.$forceUpdate();

      SetDBConf(this.DBConf);
    },

    async addTab(targetName, c) {
      if (c.Type === "" && this.DBConf[targetName]) {
        this.$message.error("duplicate connection name");
        return;
      }
      c.Port = parseInt(c.Port);
      c.Type = "kafka";
      this.DBConf[targetName] = c;
      await SetDBConf(this.DBConf);

      this.dialogVisible = false;

      this.DBConf = await DBConf();
      this.$forceUpdate();
    },

    async OpenConn(key) {
      try {
        await OpenConn(key);
      } catch (e) {
        this.$message.error("OpenConn error: " + e);
        return;
      }
      this.$emit("onOpenConn", key);
    },
  },

  async created() {},
  async mounted() {
    this.DBConf = await DBConf();
  },
};
</script>

<style>
.text {
  font-size: 14px;
}

.box-card {
  width: 350px;
  margin: 20px;
}
.grid-container {
  display: grid;
  grid-template-columns: repeat(
    auto-fit,
    minmax(320px, 1fr)
  ); /* Automatically fit columns, minimum 100px width */
  grid-gap: 10px; /* Gap between grid items */
}
</style>
